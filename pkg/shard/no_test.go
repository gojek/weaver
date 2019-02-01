package shard

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNoStrategy(t *testing.T) {
	shardConfig := json.RawMessage(`{ "backend_name": "foobar", "backend": "http://localhost" }`)

	noStrategy, err := NewNoStrategy(shardConfig)
	require.NoError(t, err, "should not have failed to parse the shard config")

	backend, err := noStrategy.Shard("whatever")
	require.NoError(t, err, "should not have failed when finding shard")

	assert.Equal(t, "http://localhost", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestNewNoStrategyFailsWhenConfigIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`[]`)

	noStrategy, err := NewNoStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Equal(t, "json: cannot unmarshal array into Go value of type shard.NoStrategyConfig", err.Error())
	assert.Nil(t, noStrategy, "should have failed to parse the shard config")
}

func TestNewNoStrategyFailsWhenBackendURLInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"backend_name": "foobar",
		"backend": "http$://google.com"
	}`)

	noStrategy, err := NewNoStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Contains(t, err.Error(), "failed to create backend:")
	assert.Nil(t, noStrategy, "should have failed to parse the shard config")
}

func TestNewNoStrategyFailsWhenConfigValueIsInvalid(t *testing.T) {
	shardConfig := json.RawMessage(`{ "backend_name": "foobar", "backend": [] }`)

	noStrategy, err := NewNoStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Contains(t, err.Error(), "cannot unmarshal array into Go struct field NoStrategyConfig.backend of type string")
	assert.Nil(t, noStrategy, "should have failed to parse the shard config")
}

func TestNoStrategyFailWhenBackendIsNotAValidURL(t *testing.T) {
	shardConfig := json.RawMessage(`{ "server": ":" }`)

	noStrategy, err := NewNoStrategy(shardConfig)
	require.Error(t, err, "should have failed to parse the shard config")

	assert.Nil(t, noStrategy)
}
