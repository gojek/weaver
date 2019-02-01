package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	etcd "github.com/coreos/etcd/client"
	"github.com/gojekfarm/weaver/internal/config"
	"github.com/gojekfarm/weaver/pkg/logger"
	"github.com/pkg/errors"
)

const (
	// ACLKeyFormat - Format for a ACL's key in a KV Store
	ACLKeyFormat = "/%s/acls/%s/acl"
)

// ACLKey - Points to a stored ACL
type ACLKey string

// GenACLKey - Generate an ACL Key given etcd's node key
func GenACLKey(key string) ACLKey {
	return ACLKey(fmt.Sprintf("%s/acl", key))
}

func GenKey(acl *ACL, pfx string) ACLKey {
	return ACLKey(fmt.Sprintf(ACLKeyFormat, pfx, acl.ID))
}

// ETCDRouteLoader - To store and modify proxy configuration
type ETCDRouteLoader struct {
	etcdClient etcd.Client
	namespace  string
}

// PutACL - Upserts a given ACL
func (routeLoader *ETCDRouteLoader) PutACL(acl *ACL) (ACLKey, error) {
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
func (routeLoader *ETCDRouteLoader) GetACL(key ACLKey) (*ACL, error) {
	res, err := etcd.NewKeysAPI(routeLoader.etcdClient).Get(context.Background(), string(key), nil)
	if err != nil {
		return nil, fmt.Errorf("fail to GET %s with %s", key, err.Error())
	}
	acl := &ACL{}
	if err := json.Unmarshal([]byte(res.Node.Value), acl); err != nil {
		return nil, err
	}

	acl.Endpoint, err = NewEndpoint(acl.EndpointConfig)
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

// NewETCDRouteLoader - Creates a new ETCDRouteLoader (routeLoader)
func NewETCDRouteLoader() (*ETCDRouteLoader, error) {
	etcdClient, err := config.NewETCDClient()
	if err != nil {
		return nil, err
	}
	return &ETCDRouteLoader{
		etcdClient: etcdClient,
		namespace:  config.ETCDKeyPrefix(),
	}, nil
}

func initEtcd(routeLoader *ETCDRouteLoader) (etcd.KeysAPI, string) {
	key := fmt.Sprintf("/%s/acls/", routeLoader.namespace)
	etc := etcd.NewKeysAPI(routeLoader.etcdClient)

	return etc, key
}

func (routeLoader *ETCDRouteLoader) WatchRoutes(ctx context.Context, upsertRouteFunc UpsertRouteFunc, deleteRouteFunc DeleteRouteFunc) {
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
			acl := &ACL{}
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

func (routeLoader *ETCDRouteLoader) BootstrapRoutes(ctx context.Context, upsertRouteFunc UpsertRouteFunc) error {
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
