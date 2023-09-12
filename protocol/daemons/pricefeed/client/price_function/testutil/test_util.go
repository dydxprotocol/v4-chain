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

func MedianErr(a []uint64) (uint64, error) {
	return uint64(0), MedianizationError
}

func CreateResponseFromJson(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}
