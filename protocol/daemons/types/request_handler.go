package types

import (
	"context"
	"net/http"
)

// RequestHandlerImpl is the struct that implements the `RequestHandler` interface.
type RequestHandlerImpl struct {
	client *http.Client
}

// RequestHandler is an interface that handles making HTTP requests.
type RequestHandler interface {
	Get(ctx context.Context, url string) (*http.Response, error)
}

// NewRequestHandlerImpl creates a new RequestHandlerImpl. It manages making HTTP requests.
func NewRequestHandlerImpl(client *http.Client) *RequestHandlerImpl {
	return &RequestHandlerImpl{
		client: client,
	}
}

// Get wraps `http.Get` which makes an HTTP GET request to a URL and returns a response.
func (r *RequestHandlerImpl) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req)
}
