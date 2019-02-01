package server

import "context"

type UpsertRouteFunc func(*ACL) error
type DeleteRouteFunc func(*ACL) error

type RouteLoader interface {
	BootstrapRoutes(context.Context, UpsertRouteFunc) error
	WatchRoutes(context.Context, UpsertRouteFunc, DeleteRouteFunc)
}
