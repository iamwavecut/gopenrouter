package gopenrouter

import (
	"net/http"

	"github.com/iamwavecut/gopenrouter/internal/apierr"
)

func decodeErrorResponse(resp *http.Response) error {
	return apierr.DecodeErrorResponse(resp)
}

func decodeAPIErrorBody(body []byte) *APIError {
	return apierr.DecodeAPIErrorBody(body)
}
