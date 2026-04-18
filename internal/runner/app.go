package runner

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

type serviceProcess struct {
	name    string
	workDir string
	command []string
	env     []string
	cmd     *exec.Cmd
}

func Run() {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	listingPort := flag.Int("listing-port", 6300, "port for listing service")
	userPort := flag.Int("user-port", 6301, "port for user service")
	publicPort := flag.Int("public-port", 7300, "port for public api")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pythonBin := filepath.Join(rootDir, ".venv", "bin", "python")
	if _, err := os.Stat(pythonBin); err != nil {
		pythonBin = "python3"
	}

	if err := ensurePythonDependency(pythonBin, "tornado"); err != nil {
		log.Fatalf("failed to start services: %v", err)
	}

	if err := ensurePortsAvailable(*listingPort, *userPort, *publicPort); err != nil {
		log.Fatalf("failed to start services: %v", err)
	}

	services := []*serviceProcess{
		{
			name:    "listing-service",
			workDir: rootDir,
			command: []string{pythonBin, "listing_service.py", fmt.Sprintf("--port=%d", *listingPort), "--debug=true"},
		},
		{
			name:    "user-service",
			workDir: filepath.Join(rootDir, "user-service"),
			command: []string{"go", "run", "./cmd/api", fmt.Sprintf("--port=%d", *userPort)},
		},
		{
			name:    "public-api",
			workDir: filepath.Join(rootDir, "public-api"),
			command: []string{
				"go", "run", "./cmd/api",
				fmt.Sprintf("--port=%d", *publicPort),
				fmt.Sprintf("--listing-service-url=http://127.0.0.1:%d", *listingPort),
				fmt.Sprintf("--user-service-url=http://127.0.0.1:%d", *userPort),
			},
		},
	}

	if err := startAll(ctx, services); err != nil {
		stopAll(services)
		log.Fatalf("failed to start services: %v", err)
	}

	log.Printf("all services are up")
	log.Printf("listing-service: http://127.0.0.1:%d", *listingPort)
	log.Printf("user-service:    http://127.0.0.1:%d", *userPort)
	log.Printf("public-api:      http://127.0.0.1:%d", *publicPort)
	log.Printf("press Ctrl+C to stop all services")

	<-ctx.Done()
	log.Printf("shutting down services")
	stopAll(services)
}

func startAll(ctx context.Context, services []*serviceProcess) error {
	for _, svc := range services {
		if err := svc.start(); err != nil {
			return fmt.Errorf("%s: %w", svc.name, err)
		}
	}

	checks := []struct {
		name string
		url  string
	}{
		{name: "listing-service", url: serviceURL(services[0], "/listings/ping")},
		{name: "user-service", url: serviceURL(services[1], "/users")},
		{name: "public-api", url: serviceURL(services[2], "/public-api/listings")},
	}

	for _, check := range checks {
		if err := waitForHTTP(ctx, check.url, 10*time.Second); err != nil {
			return fmt.Errorf("%s health check failed: %w", check.name, err)
		}
	}

	return nil
}

func (s *serviceProcess) start() error {
	if len(s.command) == 0 {
		return errors.New("empty command")
	}

	s.cmd = exec.Command(s.command[0], s.command[1:]...)
	s.cmd.Dir = s.workDir
	s.cmd.Env = append(os.Environ(), s.env...)

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := s.cmd.Start(); err != nil {
		return err
	}

	go streamOutput(s.name, stdout)
	go streamOutput(s.name, stderr)

	return nil
}

func stopAll(services []*serviceProcess) {
	var wg sync.WaitGroup
	for _, svc := range services {
		if svc == nil || svc.cmd == nil || svc.cmd.Process == nil {
			continue
		}

		wg.Add(1)
		go func(svc *serviceProcess) {
			defer wg.Done()

			_ = svc.cmd.Process.Signal(syscall.SIGTERM)

			done := make(chan struct{})
			go func() {
				_, _ = svc.cmd.Process.Wait()
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(2 * time.Second):
				_ = svc.cmd.Process.Kill()
			}
		}(svc)
	}
	wg.Wait()
}

func waitForHTTP(ctx context.Context, url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode < 500 {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}

	return fmt.Errorf("timeout waiting for %s", url)
}

func streamOutput(name string, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		log.Printf("[%s] %s", name, scanner.Text())
	}
}

func serviceURL(svc *serviceProcess, path string) string {
	for _, arg := range svc.command {
		const prefix = "--port="
		if len(arg) > len(prefix) && arg[:len(prefix)] == prefix {
			return "http://127.0.0.1:" + arg[len(prefix):] + path
		}
	}
	return ""
}

func ensurePythonDependency(pythonBin string, module string) error {
	cmd := exec.Command(pythonBin, "-c", "import "+module)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"python dependency %q is missing for %s. Activate .venv and run: python -m pip install -r python-libs.txt",
			module,
			pythonBin,
		)
	}
	return nil
}

func ensurePortsAvailable(ports ...int) error {
	for _, port := range ports {
		if err := ensurePortAvailable(port); err != nil {
			return err
		}
	}
	return nil
}

func ensurePortAvailable(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("port %d is already in use. Stop the existing process or use make run-alt", port)
	}
	_ = listener.Close()
	return nil
}
