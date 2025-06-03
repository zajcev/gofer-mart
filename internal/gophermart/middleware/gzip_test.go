package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipMiddleware_CompressResponse(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"hello world"}`))
	})

	server := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))

	// Проверяем, что тело действительно сжато
	gr, err := gzip.NewReader(res.Body)
	assert.NoError(t, err)

	unzippedBody, err := io.ReadAll(gr)
	assert.NoError(t, err)

	assert.JSONEq(t, `{"message":"hello world"}`, string(unzippedBody))
}

func TestGzipMiddleware_NoCompression(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"no compression"}`))
	})

	server := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Empty(t, res.Header.Get("Content-Encoding"))

	body, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"message":"no compression"}`, string(body))
}

func TestGzipMiddleware_DecompressRequest(t *testing.T) {
	originalBody := []byte(`{"incoming":"compressed request"}`)

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(originalBody)
	assert.NoError(t, err)
	assert.NoError(t, zw.Close())

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.JSONEq(t, `{"incoming":"compressed request"}`, string(body))

		w.WriteHeader(http.StatusOK)
	})

	server := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/", &buf)
	req.Header.Set("Content-Encoding", "gzip")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGzipMiddleware_BadCompressedRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called on invalid gzip request")
	})

	server := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not actually gzip"))
	req.Header.Set("Content-Encoding", "gzip")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGzipMiddleware_SetGzipHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"hello":"gzip"}`))
	})

	server := GzipMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
}
