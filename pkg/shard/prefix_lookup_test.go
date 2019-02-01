package shard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrefixLookupStrategy(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R_": { "timeout": 100, "backend_name": "foobar", "backend": "http://ride-service"},
		"GK_": { "backend_name": "foobar", "backend": "http://go-kilat"}
	  },
	  "prefix_splitter": "_"
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := prefixLookupStrategy.Shard("GK_123444")
	require.NoError(t, err, "should not have failed when finding shard")

	assert.Equal(t, "http://go-kilat", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestNewPrefixLookupStrategyWithDefaultPrefixSplitter(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R-": { "timeout": 100, "backend_name": "foobar", "backend": "http://ride-service"},
		"GK-": { "backend_name": "foobar", "backend": "http://go-kilat"}
	  }
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := prefixLookupStrategy.Shard("GK-123444")
	require.NoError(t, err, "should not have failed when finding shard")

	assert.Equal(t, "http://go-kilat", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestNewPrefixLookupStrategyForNoPrefix(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R-": { "timeout": 100, "backend_name": "foobar", "backend": "http://ride-service"},
		"GK-": { "backend_name": "foobar", "backend": "http://go-kilat"},
		"default": { "backend_name": "hello", "backend": "http://sm"}
	  },
	  "prefix_splitter": "-"
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := prefixLookupStrategy.Shard("123444")
	require.NoError(t, err, "should not have failed when finding shard")

	assert.Equal(t, "http://sm", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestNewPrefixLookupStrategyFailWhenTimeoutIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R-": { "timeout": "abc", "backend": "http://ride-service"},
		"GK-": { "backend": "http://go-kilat"},
		"default": { "backend_name": "hello", "backend": "http://sm"}
	  },
	  "prefix_splitter": "-"
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, prefixLookupStrategy)
}

func TestNewPrefixLookupStrategyFailWhenNoBackendGiven(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R-": { "timeout": "abc", "backend_name": "hello"},
		"GK-": { "backend_name": "mello", "backend": "http://go-kilat"},
		"default": { "backend_name": "hello", "backend": "http://sm"}
	  },
	  "prefix_splitter": "-"
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, prefixLookupStrategy)
}

func TestNewPrefixLookupStrategyFailWhenBackendIsNotString(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R-": { "backend": 123 },
		"GK-": { "backend": "http://go-kilat"},
		"default": { "backend_name": "hello", "backend": "http://sm"}
	  },
	  "prefix_splitter": "-"
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, prefixLookupStrategy)
}

func TestNewPrefixLookupStrategyFailWhenBackendIsNotAValidURL(t *testing.T) {
	shardConfig := json.RawMessage(`{
	  "backends": {
		"R-": { "backend": ":"},
		"GK-": { "backend": "http://go-kilat"},
		"default": { "backend_name": "hello", "backend": "http://sm"}
	  },
	  "prefix_splitter": "-"
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, prefixLookupStrategy)
}

func TestNewPrefixLookupStrategyFailsWhenConfigIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`[]`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Equal(t, err.Error(), "json: cannot unmarshal array into Go value of type shard.prefixLookupConfig")
	assert.Nil(t, prefixLookupStrategy, "should have failed to parse the shard config")
}

func TestNewPrefixLookupStrategyFailsWhenConfigValueIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"foo": "hello",
		"hello": []
	}`)

	prefixLookupStrategy, err := NewPrefixLookupStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Contains(t, err.Error(), "no backends specified")
	assert.Nil(t, prefixLookupStrategy, "should have failed to parse the shard config")
}
