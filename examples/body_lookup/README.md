# Backend Lookup

In this example, will deploy etcd and weaver to kubernetes and apply a simple etcd to shard between Singapore estimator and Indonesian Estimator based on key body lookup.

### Setup
```
# To kubernets cluster in local and set current context
minikube start
minikube status # verify it is up and running

# You can check dashboard by running following command
minikube dashboard

# Deploying helm components
helm init
```

### Deploying weaver

Now we have running kubernets cluster in local. Let's deploy weaver and etcd in kubernets to play with routes.

1. Clone the repo
2. On root folder of the project, run the following commands

```sh
# Connect to kubernets docker image
eval $(minikube docker-env)

# Build docker weaver image
docker build . -t weaver:stable

# Deploy weaver to kubernets
helm upgrade --debug --install proxy-cluster ./deployment/weaver --set service.type=NodePort -f ./deployment/weaver/values-env.yaml
```

We are setting service type as NodePort so that we can access it from local machine.

We have deployed weaver successfully to kubernets under release name , you can check the same in dashboard.

### Deploying simple service

Now we have to deploy simple service to kubernets and shard request using weaver.
Navigate to examples/body_lookup/ and run the following commands.

1. Build docker image for estimate service
2. Deploy docker image to 2 sharded clusters

```
# Building docker image for estimate
docker build . -t estimate:stable

# Deploying it Singapore Cluster
helm upgrade --debug --install singapore-cluster ./examples/body_lookup/estimator -f ./examples/body_lookup/estimator/values-sg.yaml

# Deploying it to Indonesian Cluster
helm upgrade --debug --install indonesia-cluster ./examples/body_lookup/estimator -f ./examples/body_lookup/estimator/values-id.yaml
```

We have a service called estimator which is sharded (Indonesian cluster, and Singapore cluster) which returns an Amount and Currency.

### Deploying Weaver ACLS

Let's deploy acl to etcd and check weaver in action.

1. Copy acls to weaver pod
2. Load acs to etcd

We have to apply acls to etcd so that we can lookup for that acl and load it. In order to apply a acl, first will copy to one of the pod
and deploy using curl request by issuing following commands.

```sh
# You can get pod name by running this command -  kubectx get pods | grep weaver | awk '{print $1}'
kubectl cp examples/body_lookup/estimate_acl.json proxy-cluster-weaver-79fb49db6f-tng8r:/go/

# Set path in etcd using curl command
curl -v etcd:2379/v2/keys/weaver/acls/estimate/acl -XPUT --data-urlencode "value@estimate_acl.json"
```

Once we set the acl in etcd, as weaver is watching for path changes continuously it just loads the acl and starts sharding requests.

### Weaver in action


Now you have wevaer which is exposed using NodePort service type. This mean you can just shard your request based on currency lookup in body as we defined in the estimate_acl.json file.

1. Get Cluster Info
2. Send request to weaver to see response from estimator

```sh
# Get cluster ip from cluster-info
kubectl cluster-info

# Using cluster ip make a curl request to weaver
curl -X POST ${CLUSTER_IP}:${NODE_PORT}/estimate -d '{"currency": "SGD"}' # This is served by singapore shard
# {"Amount": 23.23, "Currency": "SGD"}

# Getting estimate from Indonesia shard
curl -X POST ${CLUSTER_IP}:${NODE_PORT}/estimate -d '{"currency": "IDR"}' # This is served by singapore shard
# {"Amount": 81223.23, "Currency": "IDR"}
```
