package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBackend(t *testing.T) {
	serverURL := "http://localhost"

	backendOptions := BackendOptions{}
	backend, err := NewBackend("foobar", serverURL, backendOptions)
	require.NoError(t, err, "should not have failed to create new backend")

	assert.NotNil(t, backend.Handler)
	assert.Equal(t, serverURL, backend.Server.String())
}

func TestNewBackendFailsWhenURLIsInvalid(t *testing.T) {
	serverURL := ":"

	backendOptions := BackendOptions{}
	backend, err := NewBackend("foobar", serverURL, backendOptions)
	require.Error(t, err, "should have failed to create new backend")

	assert.Nil(t, backend)
}
