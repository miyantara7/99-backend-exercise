package adapter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"user-service/internal/domain/model"
)

type FileRepository struct {
	mu     sync.Mutex
	path   string
	users  []model.User
	nextID int
}

func NewFileRepository(path string) (*FileRepository, error) {
	repo := &FileRepository{
		path:   path,
		nextID: 1,
	}
	if err := repo.load(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *FileRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}

	if err := json.Unmarshal(data, &r.users); err != nil {
		return err
	}

	for _, user := range r.users {
		if user.ID >= r.nextID {
			r.nextID = user.ID + 1
		}
	}

	return nil
}

func (r *FileRepository) saveLocked() error {
	data, err := json.MarshalIndent(r.users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o644)
}

func (r *FileRepository) List(pageNum int, pageSize int) []model.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := make([]model.User, len(r.users))
	copy(items, r.users)

	sort.Slice(items, func(i, j int) bool {
		if items[i].CreatedAt == items[j].CreatedAt {
			return items[i].ID > items[j].ID
		}
		return items[i].CreatedAt > items[j].CreatedAt
	})

	offset := (pageNum - 1) * pageSize
	if offset >= len(items) {
		return []model.User{}
	}

	end := offset + pageSize
	if end > len(items) {
		end = len(items)
	}

	result := make([]model.User, end-offset)
	copy(result, items[offset:end])
	return result
}

func (r *FileRepository) GetByID(id int) (model.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.ID == id {
			return user, true
		}
	}
	return model.User{}, false
}

func (r *FileRepository) Create(name string) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UnixMicro()
	user := model.User{
		ID:        r.nextID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	r.nextID++
	r.users = append(r.users, user)

	if err := r.saveLocked(); err != nil {
		return model.User{}, err
	}
	return user, nil
}
