package server

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type mockRouteLoader struct {
	mock.Mock
}

func (mrl *mockRouteLoader) BootstrapRoutes(ctx context.Context, upsertRouteFunc UpsertRouteFunc) error {
	args := mrl.Called(ctx, upsertRouteFunc)
	return args.Error(0)
}

func (mrl *mockRouteLoader) WatchRoutes(ctx context.Context, upsertRouteFunc UpsertRouteFunc, deleteRouteFunc DeleteRouteFunc) {
	mrl.Called(ctx, upsertRouteFunc, deleteRouteFunc)
	return
}
