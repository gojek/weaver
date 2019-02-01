package server

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	etcd "github.com/coreos/etcd/client"
	"github.com/gojektech/weaver/internal/config"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Notice: This test uses time.Sleep, TODO: fix it
type ETCDRouteLoaderSuite struct {
	suite.Suite

	ng *ETCDRouteLoader
}

func (es *ETCDRouteLoaderSuite) SetupTest() {
	config.Load()
	logger.SetupLogger()

	var err error

	es.ng, err = NewETCDRouteLoader()
	assert.NoError(es.T(), err)
}

func (es *ETCDRouteLoaderSuite) TestNewETCDRouteLoader() {
	assert.NotNil(es.T(), es.ng)
}

func TestETCDRouteLoaderSuite(tst *testing.T) {
	suite.Run(tst, new(ETCDRouteLoaderSuite))
}

func (es *ETCDRouteLoaderSuite) TestPutACL() {
	aclPut := &ACL{
		ID:        "svc-01",
		Criterion: "Method(`GET`) && Path(`/ping`)",
		EndpointConfig: &EndpointConfig{
			ShardFunc: "lookup",
			Matcher:   "path",
			ShardExpr: "*",
			ShardConfig: json.RawMessage(`{
				"GF-": {
					"backend_name": "foobar",
					"backend":      "http://customer-locations-primary"
				},
				"R-": {
					"timeout":      100.0,
					"backend_name": "foobar",
					"backend":      "http://customer-locations-secondary"
				}
			}`),
		},
	}

	key, err := es.ng.PutACL(aclPut)
	assert.Nil(es.T(), err, "fail to PUT %s", aclPut)
	aclGet, err := es.ng.GetACL(key)
	assert.Nil(es.T(), err, "fail to GET with %s", key)
	assert.Equal(es.T(), aclPut.ID, aclGet.ID, "PUT %s =/= GET %s", aclPut, aclGet)
	assert.Nil(es.T(), es.ng.DelACL(key), "fail to DELETE %+v", aclPut)
}

func (es *ETCDRouteLoaderSuite) TestBootstrapRoutes() {
	aclPut := &ACL{
		ID:        "svc-01",
		Criterion: "Method(`GET`) && Path(`/ping`)",
		EndpointConfig: &EndpointConfig{
			ShardFunc:   "lookup",
			Matcher:     "path",
			ShardExpr:   "*",
			ShardConfig: json.RawMessage(`{}`),
		},
	}
	key, err := es.ng.PutACL(aclPut)
	assert.NoError(es.T(), err, "failed to PUT %s", aclPut)

	aclsChan := make(chan *ACL, 1)
	es.ng.BootstrapRoutes(context.Background(), genRouteProcessorMock(aclsChan))

	deepEqual(es.T(), aclPut, <-aclsChan)
	assert.Nil(es.T(), es.ng.DelACL(key), "fail to DELETE %+v", aclPut)
}

func (es *ETCDRouteLoaderSuite) TestBootstrapRoutesSucceedWhenARouteUpsertFails() {
	aclPut := &ACL{
		ID:        "svc-01",
		Criterion: "Method(`GET`) && Path(`/ping`)",
		EndpointConfig: &EndpointConfig{
			ShardFunc: "lookup",
			Matcher:   "path",
			ShardExpr: "*",
			ShardConfig: json.RawMessage(`{
				"GF-": {
					"backend_name": "foobar",
					"backend":      "http://customer-locations-primary"
				},
				"R-": {
					"timeout":      100.0,
					"backend_name": "foobar",
					"backend":      "http://customer-locations-secondary"
				}
			}`),
		},
	}
	key, err := es.ng.PutACL(aclPut)
	require.NoError(es.T(), err, "failed to PUT %s", aclPut)

	err = es.ng.BootstrapRoutes(context.Background(), failingUpsertRouteFunc)
	require.NoError(es.T(), err, "should not have failed to bootstrap routes")
	assert.Nil(es.T(), es.ng.DelACL(key), "fail to DELETE %+v", aclPut)
}

func (es *ETCDRouteLoaderSuite) TestBootstrapRoutesSucceedWhenARouteDoesntExist() {
	err := es.ng.BootstrapRoutes(context.Background(), successUpsertRouteFunc)
	require.NoError(es.T(), err, "should not have failed to bootstrap routes")
}

func (es *ETCDRouteLoaderSuite) TestBootstrapRoutesSucceedWhenARouteHasInvalidData() {
	aclPut := newTestACL("path")

	value := `{ "blah": "a }`
	key := "abc"
	_, err := etcd.NewKeysAPI(es.ng.etcdClient).Set(context.Background(), key, value, nil)
	require.NoError(es.T(), err, "failed to PUT %s", aclPut)

	err = es.ng.BootstrapRoutes(context.Background(), successUpsertRouteFunc)
	require.NoError(es.T(), err, "should not have failed to bootstrap routes")
	assert.Nil(es.T(), es.ng.DelACL(ACLKey(key)), "fail to DELETE %+v", aclPut)
}

func (es *ETCDRouteLoaderSuite) TestWatchRoutesUpsertRoutesWhenRoutesSet() {
	newACL := newTestACL("path")

	aclsUpserted := make(chan *ACL, 1)

	watchCtx, cancelWatch := context.WithCancel(context.Background())
	defer cancelWatch()

	go es.ng.WatchRoutes(watchCtx, genRouteProcessorMock(aclsUpserted), successUpsertRouteFunc)
	time.Sleep(100 * time.Millisecond)

	key, err := es.ng.PutACL(newACL)
	require.NoError(es.T(), err, "fail to PUT %+v", newACL)

	deepEqual(es.T(), newACL, <-aclsUpserted)
	assert.Nil(es.T(), es.ng.DelACL(key), "fail to DELETE %+v", newACL)
}

func (es *ETCDRouteLoaderSuite) TestWatchRoutesUpsertRoutesWhenRoutesUpdated() {
	newACL := newTestACL("path")
	updatedACL := newTestACL("header")

	_, err := es.ng.PutACL(newACL)
	aclsUpserted := make(chan *ACL, 1)
	watchCtx, cancelWatch := context.WithCancel(context.Background())
	defer cancelWatch()

	go es.ng.WatchRoutes(watchCtx, genRouteProcessorMock(aclsUpserted), successUpsertRouteFunc)
	time.Sleep(100 * time.Millisecond)

	key, err := es.ng.PutACL(updatedACL)
	require.NoError(es.T(), err, "fail to PUT %+v", updatedACL)

	deepEqual(es.T(), updatedACL, <-aclsUpserted)
	assert.Nil(es.T(), es.ng.DelACL(key), "fail to DELETE %+v", updatedACL)
}

func (es *ETCDRouteLoaderSuite) TestWatchRoutesDeleteRouteWhenARouteIsDeleted() {
	newACL := newTestACL("path")

	key, err := es.ng.PutACL(newACL)
	require.NoError(es.T(), err, "fail to PUT ACL %+v", newACL)

	aclsDeleted := make(chan *ACL, 1)

	watchCtx, cancelWatch := context.WithCancel(context.Background())
	defer cancelWatch()

	go es.ng.WatchRoutes(watchCtx, successUpsertRouteFunc, genRouteProcessorMock(aclsDeleted))
	time.Sleep(100 * time.Millisecond)

	err = es.ng.DelACL(key)
	require.NoError(es.T(), err, "fail to Delete %+v", newACL)

	deepEqual(es.T(), newACL, <-aclsDeleted)
}

func newTestACL(matcher string) *ACL {
	return &ACL{
		ID:        "svc-01",
		Criterion: "Method(`GET`) && Path(`/ping`)",
		EndpointConfig: &EndpointConfig{
			ShardFunc: "lookup",
			Matcher:   matcher,
			ShardExpr: "*",
			ShardConfig: json.RawMessage(`{
				"GF-": {
					"backend_name": "foobar",
					"backend":      "http://customer-locations-primary"
				},
				"R-": {
					"timeout":      100.0,
					"backend_name": "foobar",
					"backend":      "http://customer-locations-secondary"
				}
			}`),
		},
	}
}

func genRouteProcessorMock(c chan *ACL) func(*ACL) error {
	return func(acl *ACL) error {
		c <- acl
		return nil
	}
}

func deepEqual(t *testing.T, expected *ACL, actual *ACL) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Criterion, actual.Criterion)
	assertEqualJSON(t, expected.EndpointConfig.ShardConfig, actual.EndpointConfig.ShardConfig)
	assert.Equal(t, expected.EndpointConfig.ShardFunc, actual.EndpointConfig.ShardFunc)
	assert.Equal(t, expected.EndpointConfig.Matcher, actual.EndpointConfig.Matcher)
	assert.Equal(t, expected.EndpointConfig.ShardExpr, actual.EndpointConfig.ShardExpr)
}

func assertEqualJSON(t *testing.T, json1, json2 json.RawMessage) {
	var jsonVal1 interface{}
	var jsonVal2 interface{}

	err1 := json.Unmarshal(json1, &jsonVal1)
	err2 := json.Unmarshal(json2, &jsonVal2)

	assert.NoError(t, err1, "failed to parse json string")
	assert.NoError(t, err2, "failed to parse json string")
	assert.True(t, reflect.DeepEqual(jsonVal1, jsonVal2))
}

func failingUpsertRouteFunc(acl *ACL) error {
	return errors.New("error")
}

func successUpsertRouteFunc(acl *ACL) error {
	return nil
}
