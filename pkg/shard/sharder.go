package shard

import (
	"github.com/gojektech/weaver/internal/domain"
)

type Sharder interface {
	Shard(key string) (*domain.Backend, error)
}
