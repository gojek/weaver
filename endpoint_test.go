package weaver

import (
	"encoding/json"
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
	sharder := &stubSharder{}

	endpoint, err := NewEndpoint(endpointConfig, sharder)
	require.NoError(t, err, "should not fail to create an endpoint from endpointConfig")
	assert.NotNil(t, endpoint, "should create an endpoint")
	assert.Equal(t, sharder, endpoint.sharder)
}

func TestNewEndpoint_SharderIsNil(t *testing.T) {
	endpointConfig := &EndpointConfig{
		Matcher:     "path",
		ShardExpr:   "/.*",
		ShardFunc:   "lookup",
		ShardConfig: json.RawMessage(`{}`),
	}

	endpoint, err := NewEndpoint(endpointConfig, nil)
	assert.Error(t, err, "should fail to create an endpoint when sharder is nil")
	assert.Nil(t, endpoint)
}

type stubSharder struct {
}

func (stub *stubSharder) Shard(key string) (*Backend, error) {
	return nil, nil
}
