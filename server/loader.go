package server

import (
	"context"

	"github.com/gojektech/weaver"
)

type UpsertRouteFunc func(*weaver.ACL) error
type DeleteRouteFunc func(*weaver.ACL) error

type RouteLoader interface {
	BootstrapRoutes(context.Context, UpsertRouteFunc) error
	WatchRoutes(context.Context, UpsertRouteFunc, DeleteRouteFunc)
}
