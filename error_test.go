package gopenrouter

import (
	"errors"
	"fmt"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	testCases := []struct {
		name     string
		err      APIError
		expected string
	}{
		{
			name:     "error with message",
			err:      APIError{Message: "test error"},
			expected: "test error",
		},
		{
			name:     "error with empty message",
			err:      APIError{Message: ""},
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.expected {
				t.Errorf("expected error message %q, got %q", tc.expected, tc.err.Error())
			}
		})
	}
}

func TestRequestError_Error(t *testing.T) {
	originalErr := errors.New("original error")
	reqErr := &RequestError{Err: originalErr}
	expected := fmt.Sprintf("request error: %v", originalErr)
	if reqErr.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, reqErr.Error())
	}
}

func TestRequestError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	reqErr := &RequestError{Err: originalErr}
	if !errors.Is(reqErr, originalErr) {
		t.Errorf("expected unwrapped error to be the original error")
	}
}
