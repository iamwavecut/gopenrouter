package sse

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/iamwavecut/gopenrouter/internal/apierr"
)

type Event struct {
	Event string
	Data  []byte
}

type Reader struct {
	reader   *bufio.Reader
	response *http.Response
}

func NewReader(resp *http.Response) *Reader {
	return &Reader{
		reader:   bufio.NewReader(resp.Body),
		response: resp,
	}
}

func (r *Reader) RecvEvent() (Event, error) {
	var event Event
	var data [][]byte

	for {
		line, err := r.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF && len(data) > 0 {
				event.Data = bytes.Join(data, []byte("\n"))
				return event, nil
			}
			return Event{}, err
		}

		line = bytes.TrimRight(line, "\r\n")
		if len(line) == 0 {
			if len(data) == 0 {
				continue
			}
			event.Data = bytes.Join(data, []byte("\n"))
			if bytes.Equal(event.Data, []byte("[DONE]")) {
				return Event{}, io.EOF
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

func (r *Reader) Close() {
	if r == nil || r.response == nil {
		return
	}
	r.response.Body.Close()
}

func UnwrapPayload(data []byte) ([]byte, error) {
	var wrapped struct {
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Data) > 0 {
		return wrapped.Data, nil
	}

	if apiErr := apierr.DecodeAPIErrorBody(data); apiErr != nil {
		return nil, apiErr
	}

	return data, nil
}
