package shard

import (
	"encoding/json"
	"fmt"

	"github.com/gojektech/weaver"
)

func New(name string, cfg json.RawMessage) (weaver.Sharder, error) {
	newSharder, found := shardFuncTable[name]
	if !found {
		return nil, fmt.Errorf("failed to find sharder with name '%s'", name)
	}

	return newSharder(cfg)
}

type sharderGenerator func(json.RawMessage) (weaver.Sharder, error)

var shardFuncTable = map[string]sharderGenerator{
	"lookup":        NewLookupStrategy,
	"prefix-lookup": NewPrefixLookupStrategy,
	"none":          NewNoStrategy,
	"modulo":        NewModuloStrategy,
	"hashring":      NewHashRingStrategy,
	"s2":            NewS2Strategy,
}
