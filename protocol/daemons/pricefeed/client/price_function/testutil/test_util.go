package testutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

var (
	// MedianizationError is a generic error returned when the below MedianErr method is called.
	MedianizationError = errors.New("Failed to get median")
)

// MedianErr mocks the median function and returns an error.
func MedianErr(a []uint64) (uint64, error) {
	return uint64(0), MedianizationError
}

// CreateResponseFromJson creates a http response from a json string.
func CreateResponseFromJson(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}
