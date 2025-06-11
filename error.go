package gopenrouter

import "fmt"

// ErrorResponse is the response from the OpenRouter API when an error occurs.
type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

// APIError is a wrapper for the error response from the OpenRouter API.
type APIError struct {
	Message  string         `json:"message"`
	Code     any            `json:"code"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.Message
}

// RequestError provides information about generic request errors.
type RequestError struct {
	Err error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("request error: %v", e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}
