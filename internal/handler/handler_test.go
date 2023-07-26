package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OlegVankov/verbose-umbrella/internal/storage/memory"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")
	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	zr, err := gzip.NewReader(resp.Body)
	require.NoError(t, err)
	body, err := io.ReadAll(zr)
	require.NoError(t, err)
	return resp, string(body)
}

func TestRouter(t *testing.T) {
	hdl := NewHandler(memory.NewStorage(), "1234")
	ts := httptest.NewServer(NewRouter(hdl))
	defer ts.Close()

	tests := []struct {
		url    string
		method string
		want   string
		status int
	}{
		{url: "/update/gauge/Alloc/100", method: http.MethodPost, want: "", status: http.StatusOK},
		{url: "/update/gauge/Alloc/100", method: http.MethodGet, want: "", status: http.StatusMethodNotAllowed},
		{url: "/update/counter/PollCount/999", method: http.MethodPost, want: "", status: http.StatusOK},
	}

	for _, tt := range tests {
		resp, body := testRequest(t, ts, tt.method, tt.url)
		assert.Equal(t, tt.status, resp.StatusCode)
		assert.Equal(t, tt.want, body)
		resp.Body.Close()
	}
}
