package shard

import (
	"encoding/json"
	"fmt"

	"github.com/gojekfarm/weaver/internal/domain"
	"github.com/pkg/errors"
)

func NewNoStrategy(data json.RawMessage) (Sharder, error) {
	cfg := NoStrategyConfig{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	backendOptions := domain.BackendOptions{}
	backend, err := domain.NewBackend(cfg.BackendName, cfg.BackendURL, backendOptions)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to create backend: %s: %+v", err, cfg))
	}

	return &NoStrategy{
		backend: backend,
	}, nil
}

type NoStrategy struct {
	backend *domain.Backend
}

func (ns *NoStrategy) Shard(key string) (*domain.Backend, error) {
	return ns.backend, nil
}

type NoStrategyConfig struct {
	BackendDefinition `json:",inline"`
}
