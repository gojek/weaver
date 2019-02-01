package httperror

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test503Handler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)

	Err503Handler{}.ServeHTTP(w, r)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Equal(t, "{\"errors\":[{\"code\":\"weaver:service:unavailable\",\"message\":\"Something went wrong\",\"message_title\":\"Failure\",\"message_severity\":\"failure\"}]}", w.Body.String())
}
