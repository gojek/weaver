package shard

import (
	"encoding/json"

	"github.com/gojekfarm/weaver/internal/domain"
)

func NewLookupStrategy(data json.RawMessage) (Sharder, error) {
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
	backends map[string]*domain.Backend
}

func (ls *LookupStrategy) Shard(key string) (*domain.Backend, error) {
	return ls.backends[key], nil
}
