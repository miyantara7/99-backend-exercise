package dto

type ErrorResponse struct {
	Result bool     `json:"result"`
	Errors []string `json:"errors"`
}
