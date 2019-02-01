package shard

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gojektech/weaver/internal/domain"
)

const (
	defaultBackendKey     = "default"
	defaultPrefixSplitter = "-"
)

type prefixLookupConfig struct {
	PrefixSplitter string                       `json:"prefix_splitter"`
	Backends       map[string]BackendDefinition `json:"backends"`
}

func (plg prefixLookupConfig) Validate() error {
	if len(plg.Backends) == 0 {
		return errors.New("no backends specified")
	}

	return nil
}

func NewPrefixLookupStrategy(data json.RawMessage) (Sharder, error) {
	prefixLookupConfig := &prefixLookupConfig{}

	if err := json.Unmarshal(data, &prefixLookupConfig); err != nil {
		return nil, err
	}

	if err := prefixLookupConfig.Validate(); err != nil {
		return nil, err
	}

	backends, err := toBackends(prefixLookupConfig.Backends)
	if err != nil {
		return nil, err
	}

	prefixSplitter := prefixLookupConfig.PrefixSplitter
	if prefixSplitter == "" {
		prefixSplitter = defaultPrefixSplitter
	}

	return &PrefixLookupStrategy{
		backends:       backends,
		prefixSplitter: prefixSplitter,
	}, nil
}

type PrefixLookupStrategy struct {
	backends       map[string]*domain.Backend
	prefixSplitter string
}

func (pls *PrefixLookupStrategy) Shard(key string) (*domain.Backend, error) {
	prefix := strings.SplitAfter(key, pls.prefixSplitter)[0]
	if pls.backends[prefix] == nil {
		return pls.backends[defaultBackendKey], nil
	}

	return pls.backends[prefix], nil
}
