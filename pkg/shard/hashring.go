package shard

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gojekfarm/hashring"
	"github.com/gojektech/weaver"
	"github.com/pkg/errors"
)

func NewHashRingStrategy(data json.RawMessage) (weaver.Sharder, error) {
	cfg := HashRingStrategyConfig{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	hashRing, backends, err := hashringBackends(cfg)
	if err != nil {
		return nil, err
	}

	return &HashRingStrategy{
		hashRing: hashRing,
		backends: backends,
	}, nil
}

type HashRingStrategy struct {
	hashRing *hashring.HashRingCluster
	backends map[string]*weaver.Backend
}

func (rs HashRingStrategy) Shard(key string) (*weaver.Backend, error) {
	serverName := rs.hashRing.GetServer(key)
	return rs.backends[serverName], nil
}

type HashRingStrategyConfig struct {
	TotalVirtualBackends *int                         `json:"totalVirtualBackends"`
	Backends             map[string]BackendDefinition `json:"backends"`
}

func (hrCfg HashRingStrategyConfig) Validate() error {
	if hrCfg.Backends == nil || len(hrCfg.Backends) == 0 {
		return fmt.Errorf("No Shard Backends Specified Or Specified Incorrectly")
	}

	for _, backend := range hrCfg.Backends {
		if err := backend.Validate(); err != nil {
			return errors.Wrapf(err, "failed to validate backendDefinition for backend: %s", backend.BackendName)
		}
	}

	return nil
}

func hashringBackends(cfg HashRingStrategyConfig) (*hashring.HashRingCluster, map[string]*weaver.Backend, error) {
	if cfg.TotalVirtualBackends == nil || *cfg.TotalVirtualBackends < 0 {
		defaultBackends := 1000
		cfg.TotalVirtualBackends = &defaultBackends
	}

	backendDetails := map[string]*weaver.Backend{}
	hashRingCluster := hashring.NewHashRingCluster(*cfg.TotalVirtualBackends)

	virtualNodesFound := map[int]bool{}
	maxValue := 0
	rangeRegexp := regexp.MustCompile("^([\\d]+)-([\\d]+)$")
	for k, v := range cfg.Backends {
		matches := rangeRegexp.FindStringSubmatch(k)
		if len(matches) != 3 {
			return nil, nil, fmt.Errorf("Invalid range key format: %s", k)
		}
		end, _ := strconv.Atoi(matches[2])
		start, _ := strconv.Atoi(matches[1])

		if end <= start {
			return nil, nil, fmt.Errorf("Invalid range key %d-%d for backends", start, end)
		}
		for i := start; i <= end; i++ {
			if _, ok := virtualNodesFound[i]; ok {
				return nil, nil, fmt.Errorf("Overlap seen in range key %d", i)
			}
			virtualNodesFound[i] = true

			if maxValue < i {
				maxValue = i
			}
		}
		backend, err := parseBackend(v)
		if err != nil {
			return nil, nil, err
		}
		backendDetails[backend.Name] = backend
		hashRingCluster.AddServer(backend.Name, k)
	}

	if maxValue != *cfg.TotalVirtualBackends-1 {
		return nil, nil, fmt.Errorf("Shard is out of bounds Max %d found %d", *cfg.TotalVirtualBackends-1, maxValue)
	}

	for i := 0; i < *cfg.TotalVirtualBackends; i++ {
		if _, ok := virtualNodesFound[i]; !ok {
			return nil, nil, fmt.Errorf("Shard is missing coverage for %d", i)
		}
	}

	return hashRingCluster, backendDetails, nil
}
