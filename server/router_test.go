package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gojektech/weaver"
	"github.com/gojektech/weaver/pkg/shard"
	"net/http/httptest"
	"testing"

	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RouterSuite struct {
	suite.Suite

	rtr *Router
}

func (rs *RouterSuite) SetupTest() {
	logger.SetupLogger()
	routeLoader := &mockRouteLoader{}

	rs.rtr = NewRouter(routeLoader)
	require.NotNil(rs.T(), rs.rtr)
}

func TestRouterSuite(t *testing.T) {
	suite.Run(t, new(RouterSuite))
}

func (rs *RouterSuite) TestRouteNotFound() {
	req := httptest.NewRequest("GET",
		("http://" + config.ProxyServerAddress() + "/"),
		nil)
	_, err := rs.rtr.Route(req)
	assert.Error(rs.T(), err)
}

func (rs *RouterSuite) TestRouteInvalidACL() {
	req := httptest.NewRequest("GET",
		("http://" + config.ProxyServerAddress() + "/foobar"),
		nil)
	rs.rtr.UpsertRoute(
		"Method(`GET`) && Path(`/foobar`)",
		"foobar")

	_, err := rs.rtr.Route(req)
	assert.Error(rs.T(), err)
}

func (rs *RouterSuite) TestRouteReturnsACL() {
	req := httptest.NewRequest("GET",
		("http://" + config.ProxyServerAddress() + "/R-1234"),
		nil)

	// timeout is float64 because there are no integers in json
	acl := &weaver.ACL{
		ID:        "svc-01",
		Criterion: "Method(`GET`) && PathRegexp(`/(GF-|R-).*`)",
		EndpointConfig: &weaver.EndpointConfig{
			ShardConfig: json.RawMessage(`{
				"GF-": {
					"backend_name": "foobar",
					"backend":      "http://customer-locations-primary"
				},
				"R-": {
					"timeout":      100.0,
					"backend_name": "foobar",
					"backend":      "http://iamgone"
				}
			}`),
			Matcher:   "path",
			ShardExpr: "/(GF-|R-|).*",
			ShardFunc: "lookup",
		},
	}

	sharder, err := shard.New(acl.EndpointConfig.ShardFunc, acl.EndpointConfig.ShardConfig)
	require.NoError(rs.T(), err, "should not have failed to init a sharder")

	acl.Endpoint, err = weaver.NewEndpoint(acl.EndpointConfig, sharder)
	require.NoError(rs.T(), err, "should not have failed to set endpoint")

	rs.rtr.UpsertRoute(acl.Criterion, acl)

	acl, err = rs.rtr.Route(req)
	require.NoError(rs.T(), err, "should not have failed to find a route handler")

	assert.Equal(rs.T(), "svc-01", acl.ID)
}

func (rs *RouterSuite) TestBootstrapRoutesUseBootstrapRoutesOfRouteLoader() {
	ctx := context.Background()
	routeLoader := &mockRouteLoader{}

	rtr := NewRouter(routeLoader)

	routeLoader.On("BootstrapRoutes", ctx, mock.AnythingOfType("UpsertRouteFunc")).Return(nil)

	err := rtr.BootstrapRoutes(ctx)
	require.NoError(rs.T(), err, "should not have failed to bootstrap routes")

	routeLoader.AssertExpectations(rs.T())
}

func (rs *RouterSuite) TestBootstrapRoutesUseBootstrapRoutesOfRouteLoaderFail() {
	ctx := context.Background()
	routeLoader := &mockRouteLoader{}

	rtr := NewRouter(routeLoader)

	routeLoader.On("BootstrapRoutes", ctx, mock.AnythingOfType("UpsertRouteFunc")).Return(errors.New("fail"))

	err := rtr.BootstrapRoutes(ctx)
	require.Error(rs.T(), err, "should have failed to bootstrap routes")

	routeLoader.AssertExpectations(rs.T())
}

func (rs *RouterSuite) TestWatchRouteUpdatesCallsWatchRoutesOfLoader() {
	ctx := context.Background()
	routeLoader := &mockRouteLoader{}

	rtr := NewRouter(routeLoader)

	routeLoader.On("WatchRoutes", ctx, mock.AnythingOfType("UpsertRouteFunc"), mock.AnythingOfType("DeleteRouteFunc"))

	rtr.WatchRouteUpdates(ctx)

	routeLoader.AssertExpectations(rs.T())
}
