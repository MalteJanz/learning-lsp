package rpc

import (
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMsg(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	contentLength := len(content)
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", contentLength, content)
}

// returns method string + json content as bytes
func DecodeMsg(data []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return "", nil, errors.New("No separator in jsonRPC message")
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, err
	}

	var baseMsg BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMsg); err != nil {
		return "", nil, err
	}

	return baseMsg.Method, content[:contentLength], nil
}

// bufio.Scanner split function
func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return 0, nil, nil // not enough data
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, err // failed to parse content-length, return error
	}

	if len(content) < contentLength {
		return 0, nil, nil // not enough data
	}

	totalLength := len(header) + 4 + contentLength
	return totalLength, data[:totalLength], nil
}
