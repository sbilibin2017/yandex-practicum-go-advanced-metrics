package types

// APIError represents an error response returned by the API.
//
// It contains an HTTP status code and an error message intended for the client.
type APIError struct {
	Code    int    `json:"code"`    // Code is the HTTP status code associated with the error.
	Message string `json:"message"` // Message is a human-readable error message.
}
