package weaver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gojektech/weaver/pkg/matcher"
	"github.com/pkg/errors"
)

// EndpointConfig - Defines a config for external service
type EndpointConfig struct {
	Matcher     string          `json:"matcher"`
	ShardExpr   string          `json:"shard_expr"`
	ShardFunc   string          `json:"shard_func"`
	ShardConfig json.RawMessage `json:"shard_config"`
}

func (endpointConfig *EndpointConfig) genShardKeyFunc() (shardKeyFunc, error) {
	matcherFunc, found := matcher.New(endpointConfig.Matcher)
	if !found {
		return nil, errors.WithStack(fmt.Errorf("failed to find a matcherMux for: %s", endpointConfig.Matcher))
	}

	return func(req *http.Request) (string, error) {
		return matcherFunc(req, endpointConfig.ShardExpr)
	}, nil
}

type Endpoint struct {
	sharder      Sharder
	shardKeyFunc shardKeyFunc
}

func NewEndpoint(endpointConfig *EndpointConfig, sharder Sharder) (*Endpoint, error) {
	if sharder == nil {
		return nil, errors.New("nil sharder passed in")
	}

	shardKeyFunc, err := endpointConfig.genShardKeyFunc()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate shardKeyFunc for %s", endpointConfig.ShardExpr)
	}

	return &Endpoint{
		sharder:      sharder,
		shardKeyFunc: shardKeyFunc,
	}, nil
}

func (endpoint *Endpoint) Shard(request *http.Request) (*Backend, error) {
	shardKey, err := endpoint.shardKeyFunc(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find shardKey")
	}

	return endpoint.sharder.Shard(shardKey)
}

type shardKeyFunc func(*http.Request) (string, error)
