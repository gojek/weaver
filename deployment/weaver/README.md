# Deploying to Kubernetes

You can deploy to Kubernetes with the helm charts available in this repo.

### Deploying with ETCD

By default, helm install will deploy weaver with etcd. But you can disable deploying etcd if you want to reuse existing ETCD.

```sh
helm upgrade --debug --install proxy-cluster ./deployment/weaver -f ./deployment/weaver/values-env.yaml
```

This will deploy weaver with env values specified in the values-env.yaml file. In case if you want to expose weaver to outside kubernetes you can use NodePort to do that. 

```sh
helm upgrade --debug --install proxy-cluster ./deployment/weaver --set service.type=NodePort -f ./deployment/weaver/values-env.yaml
```

This will deploy along with service of type NodePort. So you can access weaver from outside your kube cluster using kube cluster address and NodePort. In production, you might want to consider ingress/load balancer.


### Deploying without ETCD

You can disable deploying ETCD in case if you want to use existing ETCD in your cluster. To do this, all you have to do is to pass `etcd.enabled` value from command line. 

```sh
helm upgrade --debug --install proxy-cluster ./deployment/weaver --set etcd.enabled=false -f ./deployment/weaver/values-env.yaml
```

This will disable deploying etcd to cluster. But you have to pass etcd host env variable `ETCD_ENDPOINTS` to make weaver work.


### Bucket List

1. Helm charts here won't support deploying statsd and sentry yet.
2. NEWRELIC key can be set to anything if you don't want to monitor your app using newrelic.
3. Similarly statsd and sentry host can be set to any value.

