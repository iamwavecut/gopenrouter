package gopenrouter

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var (
	headerData  = []byte("data: ")
	errorPrefix = []byte(`{"error":`)
)

type StreamReader struct {
	reader   *bufio.Reader
	response *http.Response
}

func (s *StreamReader) Recv() ([]byte, error) {
	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, headerData) {
			continue
		}

		line = bytes.TrimPrefix(line, headerData)
		if string(line) == "[DONE]" {
			return nil, io.EOF
		}

		if bytes.HasPrefix(line, errorPrefix) {
			return nil, fmt.Errorf("stream error: %s", line)
		}

		return line, nil
	}
}

func (s *StreamReader) Close() {
	s.response.Body.Close()
}
