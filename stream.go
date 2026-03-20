package gopenrouter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

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
