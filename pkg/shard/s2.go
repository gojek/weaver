package shard

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gojekfarm/weaver/internal/domain"
	"github.com/gojekfarm/weaver/pkg/util"
	geos2 "github.com/golang/geo/s2"
)

var (
	defaultBackendS2id = "default"
)

func NewS2Strategy(data json.RawMessage) (Sharder, error) {
	cfg := S2StrategyConfig{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	s2Backends := make(map[string]*domain.Backend, len(cfg.Backends))
	for s2id, backend := range cfg.Backends {
		var err error
		if s2Backends[s2id], err = parseBackend(backend); err != nil {
			return nil, err
		}
	}

	if cfg.ShardKeyPosition == nil {
		defaultPos := -1
		cfg.ShardKeyPosition = &defaultPos
	}

	return &S2Strategy{
		backends:          s2Backends,
		shardKeySeparator: cfg.ShardKeySeparator,
		shardKeyPosition:  *cfg.ShardKeyPosition,
	}, nil
}

type S2Strategy struct {
	backends          map[string]*domain.Backend
	shardKeySeparator string
	shardKeyPosition  int
}

type S2StrategyConfig struct {
	ShardKeySeparator string                       `json:"shard_key_separator"`
	ShardKeyPosition  *int                         `json:"shard_key_position,omitempty"`
	Backends          map[string]BackendDefinition `json:"backends"`
}

func (s2cfg S2StrategyConfig) Validate() error {
	if s2cfg.ShardKeySeparator == "" {
		return Error("missing required config: shard_key_separator")
	}
	if err := s2cfg.validateS2IDs(); err != nil {
		return err
	}
	return nil
}

func (s2cfg S2StrategyConfig) validateS2IDs() error {
	backendCount := len(s2cfg.Backends)
	s2IDs := make([]uint64, backendCount)
	for k := range s2cfg.Backends {
		if k != defaultBackendS2id {
			id, err := strconv.ParseUint(k, 10, 64)
			if err != nil {
				return fmt.Errorf("[error] Bad S2 ID found in backends: %s", k)
			}
			s2IDs = append(s2IDs, id)
		}
	}

	if util.ContainsOverlappingS2IDs(s2IDs) {
		return fmt.Errorf("[error]  Overlapping S2 IDs found in backends: %v", s2cfg.Backends)
	}
	return nil
}

func s2idFromLatLng(latLng []string) (s2id geos2.CellID, err error) {
	if len(latLng) != 2 {
		err = Error("lat lng key is not valid")
		return
	}
	lat, err := strconv.ParseFloat(latLng[0], 64)
	if err != nil {
		err = Error("fail to parse latitude")
		return
	}
	lng, err := strconv.ParseFloat(latLng[1], 64)
	if err != nil {
		err = Error("fail to parse longitude")
		return
	}
	s2LatLng := geos2.LatLngFromDegrees(lat, lng)
	if !s2LatLng.IsValid() {
		err = Error("fail to convert lat-long to geos2 objects")
		return
	}
	s2id = geos2.CellIDFromLatLng(s2LatLng)
	return
}

func s2idFromSmartID(smartIDComponents []string, pos int) (s2id geos2.CellID, err error) {
	if len(smartIDComponents) <= pos {
		err = Error("failed to get location from smart-id")
		return
	}
	s2idStr := smartIDComponents[pos]
	s2idUint, err := strconv.ParseUint(s2idStr, 10, 64)
	if err != nil {
		err = Error("failed to parse s2id")
		return
	}
	s2id = geos2.CellID(s2idUint)
	return

}

func (s2s *S2Strategy) s2ID(key string) (uint64, error) {
	shardKeyComponents := strings.Split(key, s2s.shardKeySeparator)
	var s2CellID geos2.CellID
	var err error

	switch s2s.shardKeyPosition {
	case -1:
		s2CellID, err = s2idFromLatLng(shardKeyComponents)
	default:
		s2CellID, err = s2idFromSmartID(shardKeyComponents, s2s.shardKeyPosition)
	}

	return uint64(s2CellID), err
}

func (s2s *S2Strategy) Shard(key string) (*domain.Backend, error) {
	s2id, err := s2s.s2ID(key)
	if err != nil {
		return nil, err
	}

	s2CellID := geos2.CellID(s2id)

	for s2Str, backendConfig := range s2s.backends {
		cellInt, err := strconv.ParseUint(s2Str, 10, 64)
		if err != nil {
			continue
		}
		shardCellID := geos2.CellID(cellInt)
		if shardCellID.Contains(s2CellID) {
			return backendConfig, nil
		}
	}

	if _, ok := s2s.backends[defaultBackendS2id]; ok {
		return s2s.backends[defaultBackendS2id], nil
	}

	return nil, Error("fail to find backend")
}
