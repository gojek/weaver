package shard

import (
	"encoding/json"
	"strconv"

	"github.com/gojektech/weaver/internal/domain"
	"github.com/pkg/errors"
)

func NewModuloStrategy(data json.RawMessage) (Sharder, error) {
	shardConfig := map[string]BackendDefinition{}
	if err := json.Unmarshal(data, &shardConfig); err != nil {
		return nil, err
	}

	backends, err := toBackends(shardConfig)
	if err != nil {
		return nil, err
	}

	return &ModuloStrategy{
		backends: backends,
	}, nil
}

type ModuloStrategy struct {
	backends map[string]*domain.Backend
}

func (ms ModuloStrategy) Shard(key string) (*domain.Backend, error) {
	id, err := strconv.Atoi(key)
	if err != nil {
		return nil, errors.Wrapf(err, "not an integer key: %s", key)
	}

	modulo := id % (len(ms.backends))
	return ms.backends[strconv.Itoa(modulo)], nil
}
