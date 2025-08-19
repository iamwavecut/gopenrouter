package gopenrouter

import "fmt"

// ErrorResponse is the response from the OpenRouter API when an error occurs.
type ErrorResponse struct {
	Error *APIError `json:"error,omitempty"`
}

// ProviderError provides provider-specific error details (if available).
type ProviderError map[string]any

func (e *ProviderError) Message() any {
	if e == nil {
		return nil
	}
	if msg, ok := (*e)["message"]; ok {
		return msg
	}
	return nil
}

// APIError is a wrapper for the error response from the OpenRouter API.
type APIError struct {
	Message       string         `json:"message"`
	Code          any            `json:"code"`
	Metadata      map[string]any `json:"metadata,omitempty"`
	ProviderError *ProviderError `json:"provider_error,omitempty"`
}

// Error returns the error message.
func (e *APIError) Error() string {
	return e.Message
}

// RequestError provides information about generic request errors.
type RequestError struct {
	HTTPStatus     string
	HTTPStatusCode int
	Err            error
	Body           []byte
}

func (e *RequestError) Error() string {
	if e == nil {
		return "request error: <nil>"
	}
	if e.HTTPStatus != "" {
		return fmt.Sprintf("request error: %s: %v", e.HTTPStatus, e.Err)
	}
	return fmt.Sprintf("request error: %v", e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}
