package types

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
