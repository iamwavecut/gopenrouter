package gopenrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) newRequest(ctx context.Context, method, rawURL string, payload any) (*http.Request, error) {
	return c.newRequestWithQuery(ctx, method, rawURL, nil, payload)
}

func (c *Client) newRequestWithQuery(ctx context.Context, method, rawURL string, query url.Values, payload any) (*http.Request, error) {
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return nil, &RequestError{Err: err}
		}
	}

	if len(query) > 0 {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return nil, &RequestError{Err: err}
		}
		values := parsed.Query()
		for key, items := range query {
			for _, item := range items {
				values.Add(key, item)
			}
		}
		parsed.RawQuery = values.Encode()
		rawURL = parsed.String()
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, &body)
	if err != nil {
		return nil, &RequestError{Err: err}
	}

	req.Header.Set("Authorization", "Bearer "+c.config.authToken())
	if c.config.SiteURL != "" {
		req.Header.Set("HTTP-Referer", c.config.SiteURL)
	}
	if c.config.SiteName != "" {
		req.Header.Set("X-OpenRouter-Title", c.config.SiteName)
		req.Header.Set("X-Title", c.config.SiteName)
	}
	if len(c.config.SiteCategories) > 0 {
		req.Header.Set("X-OpenRouter-Categories", strings.Join(c.config.SiteCategories, ","))
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) doJSON(ctx context.Context, method, rawURL string, query url.Values, payload any, out any) error {
	req, err := c.newRequestWithQuery(ctx, method, rawURL, query, payload)
	if err != nil {
		return err
	}

	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return &RequestError{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return decodeErrorResponse(resp)
	}
	if out == nil {
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			return &RequestError{
				HTTPStatus:     resp.Status,
				HTTPStatusCode: resp.StatusCode,
				Err:            fmt.Errorf("failed to discard response body: %w", err),
			}
		}
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return &RequestError{
			HTTPStatus:     resp.Status,
			HTTPStatusCode: resp.StatusCode,
			Err:            fmt.Errorf("failed to decode response: %w", err),
		}
	}

	return nil
}

func (c *ClientConfig) authToken() string {
	if c.AuthToken != "" {
		return c.AuthToken
	}
	return c.APIKey
}
