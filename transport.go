package gopenrouter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/iamwavecut/gopenrouter/internal/apierr"
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

func decodeErrorResponse(resp *http.Response) error {
	return apierr.DecodeErrorResponse(resp)
}

func decodeAPIErrorBody(body []byte) *APIError {
	return apierr.DecodeAPIErrorBody(body)
}

type SSEEvent struct {
	Event string
	Data  []byte
}

type SSEReader struct {
	reader   *bufio.Reader
	response *http.Response
}

func (s *SSEReader) RecvEvent() (SSEEvent, error) {
	var event SSEEvent
	var data [][]byte

	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF && len(data) > 0 {
				event.Data = bytes.Join(data, []byte("\n"))
				return event, nil
			}
			return SSEEvent{}, err
		}

		line = bytes.TrimRight(line, "\r\n")
		if len(line) == 0 {
			if len(data) == 0 {
				continue
			}
			event.Data = bytes.Join(data, []byte("\n"))
			if bytes.Equal(event.Data, []byte("[DONE]")) {
				return SSEEvent{}, io.EOF
			}
			return event, nil
		}
		if bytes.HasPrefix(line, []byte(":")) {
			continue
		}

		field, value, ok := bytes.Cut(line, []byte(":"))
		if !ok {
			continue
		}
		value = bytes.TrimLeft(value, " ")

		switch string(field) {
		case "event":
			event.Event = string(value)
		case "data":
			data = append(data, value)
		}
	}
}

type StreamReader struct {
	sse *SSEReader
}

func (s *StreamReader) Recv() ([]byte, error) {
	event, err := s.sse.RecvEvent()
	if err != nil {
		return nil, err
	}
	return unwrapSSEPayload(event.Data)
}

func (s *StreamReader) Close() {
	if s == nil || s.sse == nil || s.sse.response == nil {
		return
	}
	s.sse.response.Body.Close()
}

func newStreamReader(resp *http.Response) *StreamReader {
	return &StreamReader{
		sse: &SSEReader{
			reader:   bufio.NewReader(resp.Body),
			response: resp,
		},
	}
}

func unwrapSSEPayload(data []byte) ([]byte, error) {
	var wrapped struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Data) > 0 {
		return wrapped.Data, nil
	}

	if apiErr := decodeAPIErrorBody(data); apiErr != nil {
		return nil, apiErr
	}

	return data, nil
}
