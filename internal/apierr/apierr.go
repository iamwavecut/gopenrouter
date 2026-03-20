package apierr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iamwavecut/gopenrouter/shared"
)

func DecodeErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &shared.RequestError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			Err:            fmt.Errorf("failed to read error response: %w", err),
		}
	}

	apiErr := DecodeAPIErrorBody(body)
	if apiErr != nil {
		return apiErr
	}

	return &shared.RequestError{
		HTTPStatus:     resp.Status,
		HTTPStatusCode: resp.StatusCode,
		Err:            fmt.Errorf("request failed"),
		Body:           body,
	}
}

func DecodeAPIErrorBody(body []byte) *shared.APIError {
	if len(body) == 0 {
		return nil
	}

	var openRouter shared.ErrorResponse
	if err := json.Unmarshal(body, &openRouter); err == nil {
		if openRouter.Error != nil && openRouter.Error.Message != "" {
			if openRouter.Error.Type == "" && openRouter.Type != "" {
				openRouter.Error.Type = openRouter.Type
			}
			return openRouter.Error
		}
	}

	var anthropic struct {
		Type  string           `json:"type"`
		Error *shared.APIError `json:"error"`
	}
	if err := json.Unmarshal(body, &anthropic); err == nil && anthropic.Error != nil && anthropic.Error.Message != "" {
		if anthropic.Error.Type == "" {
			anthropic.Error.Type = anthropic.Type
		}
		return anthropic.Error
	}

	return nil
}
