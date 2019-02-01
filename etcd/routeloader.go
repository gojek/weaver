package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gojektech/weaver"
	"github.com/gojektech/weaver/pkg/shard"
	"sort"

	etcd "github.com/coreos/etcd/client"
	"github.com/gojektech/weaver/config"
	"github.com/gojektech/weaver/pkg/logger"
	"github.com/gojektech/weaver/server"
	"github.com/pkg/errors"
)

func NewRouteLoader() (*ETCDRouteLoader, error) {
	etcdClient, err := config.NewETCDClient()
	if err != nil {
		return nil, err
	}
	return &ETCDRouteLoader{
		etcdClient: etcdClient,
		namespace:  config.ETCDKeyPrefix(),
	}, nil
}

// ETCDRouteLoader - To store and modify proxy configuration
type ETCDRouteLoader struct {
	etcdClient etcd.Client
	namespace  string
}

// PutACL - Upserts a given ACL
func (routeLoader *ETCDRouteLoader) PutACL(acl *server.ACL) (ACLKey, error) {
	key := GenKey(acl, routeLoader.namespace)
	val, err := json.Marshal(acl)
	if err != nil {
		return "", err
	}
	_, err = etcd.NewKeysAPI(routeLoader.etcdClient).Set(context.Background(), string(key), string(val), nil)
	if err != nil {
		return "", fmt.Errorf("fail to PUT %s:%s with %s", key, acl, err.Error())
	}
	return key, nil
}

// GetACL - Fetches an ACL given an ACLKey
func (routeLoader *ETCDRouteLoader) GetACL(key ACLKey) (*server.ACL, error) {
	res, err := etcd.NewKeysAPI(routeLoader.etcdClient).Get(context.Background(), string(key), nil)
	if err != nil {
		return nil, fmt.Errorf("fail to GET %s with %s", key, err.Error())
	}
	acl := &server.ACL{}
	if err := json.Unmarshal([]byte(res.Node.Value), acl); err != nil {
		return nil, err
	}

	sharder, err := shard.New(acl.EndpointConfig.ShardFunc, acl.EndpointConfig.ShardConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialize sharder '%s'", acl.EndpointConfig.ShardFunc)
	}

	acl.Endpoint, err = weaver.NewEndpoint(acl.EndpointConfig, sharder)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a new Endpoint for key: %s", key)
	}

	return acl, nil
}

// DelACL - Deletes an ACL given an ACLKey
func (routeLoader *ETCDRouteLoader) DelACL(key ACLKey) error {
	_, err := etcd.NewKeysAPI(routeLoader.etcdClient).Delete(context.Background(), string(key), nil)
	if err != nil {
		return fmt.Errorf("fail to DELETE %s with %s", key, err.Error())
	}
	return nil
}

func (routeLoader *ETCDRouteLoader) WatchRoutes(ctx context.Context, upsertRouteFunc server.UpsertRouteFunc, deleteRouteFunc server.DeleteRouteFunc) {
	etc, key := initEtcd(routeLoader)
	watcher := etc.Watcher(key, &etcd.WatcherOptions{Recursive: true})

	logger.Infof("starting etcd watcher on %s", key)
	for {
		res, err := watcher.Next(ctx)
		if err != nil {
			logger.Errorf("stopping etcd watcher on %s: %v", key, err)
			return
		}

		logger.Debugf("registered etcd watcher event on %v with action %s", res, res.Action)
		switch res.Action {
		case "set":
			fallthrough
		case "update":
			logger.Debugf("fetching node key %s", res.Node.Key)
			acl, err := routeLoader.GetACL(ACLKey(res.Node.Key))
			if err != nil {
				logger.Errorf("error in fetching %s: %v", res.Node.Key, err)
				continue
			}

			logger.Infof("upserting %v to router", acl)
			err = upsertRouteFunc(acl)
			if err != nil {
				logger.Errorf("error in upserting %v: %v ", acl, err)
				continue
			}
		case "delete":
			acl := &server.ACL{}
			err := acl.GenACL(res.PrevNode.Value)
			if err != nil {
				logger.Errorf("error in unmarshalling %s: %v", res.PrevNode.Value, err)
				continue
			}

			logger.Infof("deleteing %v to router", acl)
			err = deleteRouteFunc(acl)
			if err != nil {
				logger.Errorf("error in deleting %v: %v ", acl, err)
				continue
			}
		}
	}
}

func (routeLoader *ETCDRouteLoader) BootstrapRoutes(ctx context.Context, upsertRouteFunc server.UpsertRouteFunc) error {
	// TODO: Consider error scenarios and return an error when it makes sense
	etc, key := initEtcd(routeLoader)
	logger.Infof("bootstrapping router using etcd on %s", key)
	res, err := etc.Get(ctx, key, nil)
	if err != nil {
		logger.Infof("creating router namespace on etcd using %s", key)
		_, _ = etc.Set(ctx, key, "", &etcd.SetOptions{
			Dir: true,
		})
	}

	if res != nil {
		sort.Sort(res.Node.Nodes)
		for _, nd := range res.Node.Nodes {
			logger.Debugf("fetching node key %s", nd.Key)
			acl, err := routeLoader.GetACL(GenACLKey(nd.Key))
			if err != nil {
				logger.Errorf("error in fetching %s: %v", nd.Key, err)
				continue
			}

			logger.Infof("upserting %v to router", acl)
			err = upsertRouteFunc(acl)
			if err != nil {
				logger.Errorf("error in upserting %v: %v ", acl, err)
				continue
			}
		}
	}

	return nil
}

func initEtcd(routeLoader *ETCDRouteLoader) (etcd.KeysAPI, string) {
	key := fmt.Sprintf("/%s/acls/", routeLoader.namespace)
	etc := etcd.NewKeysAPI(routeLoader.etcdClient)

	return etc, key
}
