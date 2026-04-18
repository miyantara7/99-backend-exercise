package userservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"public-api/internal/domain/dto"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (c *Client) CreateUser(name string) (dto.UserResult, error) {
	form := url.Values{}
	form.Set("name", name)

	resp, err := c.httpClient.Post(
		c.baseURL+"/users",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return dto.UserResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.UserResult{}, err
	}

	var payload dto.ClientResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return dto.UserResult{}, err
	}

	return dto.UserResult{
		Result: dto.Result{Status: resp.StatusCode, Errors: payload.Errors},
		User:   payload.User,
	}, nil
}

func (c *Client) GetUser(userID int) (dto.UserResult, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/users/"+strconv.Itoa(userID), nil)
	if err != nil {
		return dto.UserResult{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return dto.UserResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto.UserResult{}, err
	}

	var payload dto.ClientResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return dto.UserResult{}, fmt.Errorf("invalid user service response: %w", err)
	}

	return dto.UserResult{
		Result: dto.Result{Status: resp.StatusCode, Errors: payload.Errors},
		User:   payload.User,
	}, nil
}
