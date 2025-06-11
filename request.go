package gopenrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (c *Client) newRequest(ctx context.Context, method, url string, payload any) (*http.Request, error) {
	var body *bytes.Buffer
	if payload != nil {
		body = &bytes.Buffer{}
		if err := json.NewEncoder(body).Encode(payload); err != nil {
			return nil, &RequestError{Err: err}
		}
	} else {
		body = &bytes.Buffer{}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, &RequestError{Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+c.config.AuthToken)
	if c.config.SiteURL != "" {
		req.Header.Set("HTTP-Referer", c.config.SiteURL)
	}
	if c.config.SiteName != "" {
		req.Header.Set("X-Title", c.config.SiteName)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
