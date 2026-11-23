package ws

import (
	"net/http"
	"net/http/httptest"
	_ "net/http/pprof" // Registers on http.DefaultServeMux, as the cometbft fork does.
	"testing"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/require"
)

func TestMuxServesOnlyWs(t *testing.T) {
	defaultMux := httptest.NewServer(http.DefaultServeMux)
	defer defaultMux.Close()
	require.Equal(t, http.StatusOK, get(t, defaultMux.URL+"/debug/pprof/"))

	srv := httptest.NewServer(NewWebsocketServer(nil, nil, log.NewNopLogger(), 0).mux())
	defer srv.Close()
	require.Equal(t, http.StatusNotFound, get(t, srv.URL+"/debug/pprof/"))
	require.Equal(t, http.StatusNotFound, get(t, srv.URL+"/debug/pprof/heap"))
	require.Equal(t, http.StatusNotFound, get(t, srv.URL+"/"))
	// Not a websocket upgrade request.
	require.Equal(t, http.StatusBadRequest, get(t, srv.URL+"/ws"))
}

func get(t *testing.T, url string) int {
	resp, err := http.Get(url)
	require.NoError(t, err)
	resp.Body.Close()
	return resp.StatusCode
}
