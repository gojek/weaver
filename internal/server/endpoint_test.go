package server

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEndpoint(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:     "path",
		ShardExpr:   "/.*",
		ShardFunc:   "lookup",
		ShardConfig: json.RawMessage(`{}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.NoError(t, err, "should not fail to create an endpoint from endpointConfig")
	assert.NotNil(t, endpoint, "should create an endpoint")
}

func TestNewEndpointWhenShardFuncIsMissing(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:     "path",
		ShardExpr:   "/.*",
		ShardFunc:   "invalid",
		ShardConfig: json.RawMessage(`{}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.Error(t, err, "should fail to create an endpoint when ShardFunc is missing")
	assert.Equal(t, "failed to find ShardFunc for: invalid", err.Error())

	assert.Nil(t, endpoint, "should fail to create an endpoint when ShardFunc is missing")
}

func TestNewEndpointWhenBackendNameIsMissing(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:   "path",
		ShardExpr: "/.*",
		ShardFunc: "lookup",
		ShardConfig: json.RawMessage(`{
			"R-": {
				"timeout": 100.0,
				"backend": "http://iamgone"
			}
		}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.Error(t, err, "should fail to create an endpoint when backend name is missing")
	assert.Contains(t, err.Error(), "failed to get sharder for /.*: failed to validate backend definition: missing backend name in shard config:")
	assert.Nil(t, endpoint, "should fail to create an endpoint when backend name is missing")
}

func TestNewEndpointWhenMatcherIsMissing(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:     "invalid",
		ShardExpr:   "/.*",
		ShardFunc:   "lookup",
		ShardConfig: json.RawMessage(`{}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.Error(t, err, "should fail to create an endpoint when ShardFunc is missing")
	assert.Equal(t, "failed to generate shardKeyFunc for /.*: failed to find a matcherMux for: invalid", err.Error())

	assert.Nil(t, endpoint, "should fail to create an endpoint when ShardFunc is missing")
}

func TestShardEndpoint(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:   "path",
		ShardExpr: "/(.*)",
		ShardFunc: "lookup",
		ShardConfig: json.RawMessage(`{
			"GK-": {
				"backend_name": "foobar",
				"backend":      "http://localhost"
			},
			"R-": {
				"backend_name": "foobar",
				"backend":      "http://localhost/r"
			}
		}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.NoError(t, err, "should fail to create an endpoint when Sharder cannot be created")

	req := httptest.NewRequest("GET", "/R-", nil)
	backend, err := endpoint.Shard(req)
	require.NoError(t, err, "should not have failed to shard endpoint")

	assert.Equal(t, "http://localhost/r", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestShardEndpointModuloSuccess(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:   "path",
		ShardExpr: "/(.*)",
		ShardFunc: "modulo",
		ShardConfig: json.RawMessage(`{
			"0": {
				"backend_name": "foobar",
				"backend":      "http://localhost"
			},
			"1": {
				"backend_name": "foobar",
				"backend":      "http://localhost/r"
			}
		}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.NoError(t, err, "should fail to create an endpoint when Sharder cannot be created")

	req := httptest.NewRequest("GET", "/1234", nil)
	backend, err := endpoint.Shard(req)
	require.NoError(t, err, "should not have failed to shard endpoint")

	assert.Equal(t, "http://localhost", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}

func TestShardEndpointModuloFail(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:   "path",
		ShardExpr: "/(.*)",
		ShardFunc: "modulo",
		ShardConfig: json.RawMessage(`{
			"0": {
				"backend_name": "foobar",
				"backend":      "http://localhost"
			},
			"1": {
				"backend_name": "foobar",
				"backend":      "http://localhost/r"
			}
		}`),
	}

	endpoint, err := NewEndpoint(endpointConfig)
	require.NoError(t, err, "should fail to create an endpoint when Sharder cannot be created")

	req := httptest.NewRequest("GET", "/abc", nil)
	backend, err := endpoint.Shard(req)
	require.Error(t, err, "should have failed to shard endpoint")

	assert.Nil(t, backend)
}

func TestShardEndpointHashringSuccess(t *testing.T) {
	shardConfig := json.RawMessage(`{
		"totalVirtualBackends": 1000,
		"backends": {
			"0-250": { "timeout": 100, "backend_name": "foobar1", "backend": "http://shard00.local"},
			"251-500": { "backend_name": "foobar2", "backend": "http://shard01.local"},
			"501-725": { "backend_name": "foobar3", "backend": "http://shard02.local"},
			"726-999": { "backend_name": "foobar4", "backend": "http://shard03.local"}
		}
	}`)
	endpointConfig := &EndpointConfig{
		Matcher:     "path",
		ShardExpr:   "/(.*)",
		ShardFunc:   "hashring",
		ShardConfig: shardConfig,
	}
	endpoint, err := NewEndpoint(endpointConfig)
	require.NoError(t, err, "should fail to create an endpoint when Sharder cannot be created")

	req := httptest.NewRequest("GET", "/1234", nil)
	backend, err := endpoint.Shard(req)
	require.NoError(t, err, "should not have failed to shard endpoint")
	assert.Equal(t, "http://shard01.local", backend.Server.String())
	assert.NotNil(t, backend.Handler)

	req = httptest.NewRequest("GET", "/abc", nil)
	backend, err = endpoint.Shard(req)
	require.NoError(t, err, "should not have failed to shard endpoint")
	assert.Equal(t, "http://shard02.local", backend.Server.String())
	assert.NotNil(t, backend.Handler)
}
