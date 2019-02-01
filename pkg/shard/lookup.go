package shard

import (
	"encoding/json"

	"github.com/gojektech/weaver"
)

func NewLookupStrategy(data json.RawMessage) (weaver.Sharder, error) {
	shardConfig := map[string]BackendDefinition{}
	if err := json.Unmarshal(data, &shardConfig); err != nil {
		return nil, err
	}

	backends, err := toBackends(shardConfig)
	if err != nil {
		return nil, err
	}

	return &LookupStrategy{
		backends: backends,
	}, nil
}

type LookupStrategy struct {
	backends map[string]*weaver.Backend
}

func (ls *LookupStrategy) Shard(key string) (*weaver.Backend, error) {
	return ls.backends[key], nil
}
