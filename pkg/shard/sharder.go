package shard

import (
	"github.com/gojektech/weaver"
)

type Sharder interface {
	Shard(key string) (*weaver.Backend, error)
}
