package testutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var (
	MedianizationError = errors.New("Failed to get median")
)

func CreateResponseFromJson(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}
