package shard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModuloShardStrategy(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"0": { "timeout": 100, "backend_name": "foobar", "backend": "http://shard00.local"},
		"1": { "backend_name": "foobar1", "backend": "http://shard01.local"},
		"2": { "backend_name": "foobar2", "backend": "http://shard02.local"}
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := moduloStrategy.Shard("5678987")
	require.NoError(t, err, "should not have failed to find backend")

	assert.Equal(t, "http://shard02.local", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestNewModuloShardStrategyFailure(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"0": { "timeout": 100, "backend_name": "foobar", "backend": "http://shard00.local"},
		"1": { "backend_name": "foobar1", "backend": "http://shard01.local"},
		"2": { "backend_name": "foobar2", "backend": "http://shard02.local"}
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := moduloStrategy.Shard("abcd")
	require.Error(t, err, "should have failed to find backend")

	assert.Nil(t, backend)
}

func TestNewModuloStrategyFailWhenTimeoutIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"0": { "backend_name": "A", "timeout": "abc", "backend": "http://shard00.local"},
		"1": { "backend_name": "B", "backend": "http://shard01.local"}
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, moduloStrategy)
	assert.Contains(t, err.Error(), "cannot unmarshal string into Go struct field BackendDefinition.timeout of type float64")
}

func TestNewModuloStrategyFailWhenNoBackendGiven(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"0": { "backend_name": "hello"},
		"1": { "backend_name": "mello", "backend": "http://shard01.local"}
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, moduloStrategy)
	assert.Contains(t, err.Error(), "missing backend url in shard config:")
}

func TestNewModuloStrategyFailWhenBackendIsNotString(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"0": { "backend_name": "hello", "backend": 123 },
		"1": { "backend_name": "mello", "backend": "http://shard01.local"}
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, moduloStrategy)
	assert.Contains(t, err.Error(), "cannot unmarshal number")
}

func TestNewModuloStrategyFailWhenBackendIsNotAValidURL(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"0": { "backend_name": "hello", "backend": ":"},
		"1": { "backend_name": "mello", "backend": "http://shard01.local"}
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, moduloStrategy)
	assert.Contains(t, err.Error(), "URL Parsing failed for")
}

func TestNewModuloStrategyFailsWhenConfigIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`[]`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Contains(t, err.Error(), "json: cannot unmarshal array into Go value of type map[string]shard.BackendDefinition")
	assert.Nil(t, moduloStrategy, "should have failed to parse the shard config")
}

func TestNewModuloStrategyFailsWhenConfigValueIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"foo": "hello",
		"hello": []
	}`)

	moduloStrategy, err := NewModuloStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Contains(t, err.Error(), "cannot unmarshal string into Go value of type shard.BackendDefinition")
	assert.Nil(t, moduloStrategy, "should have failed to parse the shard config")
}
