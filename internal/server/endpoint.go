package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gojekfarm/weaver/internal/domain"
	"github.com/gojekfarm/weaver/pkg/shard"
	"github.com/pkg/errors"
)

var shardFuncTable = map[string]SharderGenerator{
	"lookup":        shard.NewLookupStrategy,
	"prefix-lookup": shard.NewPrefixLookupStrategy,
	"none":          shard.NewNoStrategy,
	"modulo":        shard.NewModuloStrategy,
	"hashring":      shard.NewHashRingStrategy,
	"s2":            shard.NewS2Strategy,
}

// EndpointConfig - Defines a config for external service
type EndpointConfig struct {
	Matcher     string          `json:"matcher"`
	ShardExpr   string          `json:"shard_expr"`
	ShardFunc   string          `json:"shard_func"`
	ShardConfig json.RawMessage `json:"shard_config"`
}

type shardKeyFunc func(*http.Request) (string, error)

func (endpointConfig *EndpointConfig) genShardKeyFunc() (shardKeyFunc, error) {
	matcherFunc, found := matcherMux[endpointConfig.Matcher]
	if !found {
		return nil, errors.WithStack(fmt.Errorf("failed to find a matcherMux for: %s", endpointConfig.Matcher))
	}

	return func(req *http.Request) (string, error) {
		return matcherFunc(req, endpointConfig.ShardExpr)
	}, nil
}

type SharderGenerator func(json.RawMessage) (shard.Sharder, error)

type Endpoint struct {
	sharder      shard.Sharder
	shardKeyFunc shardKeyFunc
}

func NewEndpoint(endpointConfig *EndpointConfig) (*Endpoint, error) {
	shardFunc, found := shardFuncTable[endpointConfig.ShardFunc]
	if !found {
		return nil, errors.WithStack(fmt.Errorf("failed to find ShardFunc for: %s", endpointConfig.ShardFunc))
	}

	sharder, err := shardFunc(endpointConfig.ShardConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get sharder for %s", endpointConfig.ShardExpr)
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

func (endpoint *Endpoint) Shard(request *http.Request) (*domain.Backend, error) {
	shardKey, err := endpoint.shardKeyFunc(request)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find shardKey")
	}

	return endpoint.sharder.Shard(shardKey)
}
