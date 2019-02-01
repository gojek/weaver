package shard

import (
	"github.com/gojekfarm/weaver/internal/domain"
)

type Sharder interface {
	Shard(key string) (*domain.Backend, error)
}
