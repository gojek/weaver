package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vulcand/route"
)

type Router struct {
	route.Router
	loader RouteLoader
}

type apiName string

func (router *Router) Route(req *http.Request) (*ACL, error) {
	rt, err := router.Router.Route(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find route with url: %s", req.URL)
	}

	if rt == nil {
		return nil, errors.WithStack(fmt.Errorf("route not found: %s", req.URL))
	}

	acl, ok := rt.(*ACL)
	if !ok {
		return nil, errors.WithStack(fmt.Errorf("error in casting %v to acl", rt))
	}

	return acl, nil
}

func NewRouter(loader RouteLoader) *Router {
	return &Router{
		Router: route.New(),
		loader: loader,
	}
}

func (router *Router) WatchRouteUpdates(routeSyncCtx context.Context) {
	router.loader.WatchRoutes(routeSyncCtx, router.upsertACL, router.deleteACL)
}

func (router *Router) BootstrapRoutes(ctx context.Context) error {
	return router.loader.BootstrapRoutes(ctx, router.upsertACL)
}

func (router *Router) upsertACL(acl *ACL) error {
	return router.UpsertRoute(acl.Criterion, acl)
}

func (router *Router) deleteACL(acl *ACL) error {
	return router.RemoveRoute(acl.Criterion)
}
