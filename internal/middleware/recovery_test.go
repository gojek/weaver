package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gojekfarm/weaver/internal/config"
	"github.com/gojekfarm/weaver/pkg/logger"

	"github.com/stretchr/testify/assert"
)

type testHandler struct{}

func (th testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("failed")
}

func TestRecoverMiddleware(t *testing.T) {
	config.Load()
	logger.SetupLogger()

	r := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	Recover(testHandler{}).ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"errors\":[{\"code\":\"weaver:service:unavailable\",\"message\":\"Something went wrong\",\"message_title\":\"Failure\",\"message_severity\":\"failure\"}]}", w.Body.String())
}
