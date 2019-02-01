package httperror

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test404Handler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/hello", nil)

	NotFoundHandler(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"errors\":[{\"code\":\"weaver:route:not_found\",\"message\":\"Something went wrong\",\"message_title\":\"Failure\",\"message_severity\":\"failure\"}]}", w.Body.String())
}
