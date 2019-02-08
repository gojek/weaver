package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test404Handler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)

	notFoundError(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"errors\":[{\"code\":\"weaver:route:not_found\",\"message\":\"Something went wrong\",\"message_title\":\"Failure\",\"message_severity\":\"failure\"}]}", w.Body.String())
}

func Test500Handler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)

	internalServerError(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "{\"errors\":[{\"code\":\"weaver:service:unavailable\",\"message\":\"Something went wrong\",\"message_title\":\"Internal error\",\"message_severity\":\"failure\"}]}", w.Body.String())
}

func Test503Handler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)

	err503Handler{}.ServeHTTP(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "{\"errors\":[{\"code\":\"weaver:service:unavailable\",\"message\":\"Something went wrong\",\"message_title\":\"Failure\",\"message_severity\":\"failure\"}]}", w.Body.String())
}
