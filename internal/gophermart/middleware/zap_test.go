package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
)

func TestZapMiddleware(t *testing.T) {
	core := zapcore.NewNopCore()
	testLogger := zap.New(core)
	_ = zap.ReplaceGlobals(testLogger)

	handlerCalled := false

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusTeapot) // 418 I'm a teapot
		w.Write([]byte("short and stout"))
	})

	middleware := ZapMiddleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	assert.True(t, handlerCalled, "handler should be called")

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusTeapot, res.StatusCode)
}
