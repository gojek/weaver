package shard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLookupStrategy(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"R-": { "timeout": 100, "backend_name": "foobar", "backend": "http://ride-service"},
		"GK-": { "backend_name": "foobar", "backend": "http://go-kilat"}
	}`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := lookupStrategy.Shard("GK-")
	require.NoError(t, err, "should not have failed when finding shard")

	assert.Equal(t, "http://go-kilat", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestNewLookupStrategyFailWhenTimeoutIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"R-": { "timeout": "abc", "backend": "http://ride-service"},
		"GK-": { "backend": "http://go-kilat"}
	}`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, lookupStrategy)
}

func TestNewLookupStrategyFailWhenNoBackendGiven(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"R-": { "timeout": "abc", "backend_name": "hello"},
		"GK-": { "backend_name": "mello", "backend": "http://go-kilat"}
	}`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, lookupStrategy)
}

func TestNewLookupStrategyFailWhenBackendIsNotString(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"R-": { "backend": 123 },
		"GK-": { "backend": "http://go-kilat"}
	}`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, lookupStrategy)
}

func TestNewLookupStrategyFailWhenBackendIsNotAValidURL(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"R-": { "backend": ":"},
		"GK-": { "backend": "http://go-kilat"}
	}`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, lookupStrategy)
}

func TestNewLookupStrategyFailsWhenConfigIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`[]`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Equal(t, err.Error(), "json: cannot unmarshal array into Go value of type map[string]shard.BackendDefinition")
	assert.Nil(t, lookupStrategy, "should have failed to parse the shard config")
}

func TestNewLookupStrategyFailsWhenConfigValueIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"foo": "hello",
		"hello": []
	}`)

	lookupStrategy, err := NewLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Contains(t, err.Error(), "cannot unmarshal string into Go value of type shard.BackendDefinition")
	assert.Nil(t, lookupStrategy, "should have failed to parse the shard config")
}
