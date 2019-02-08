## Weaver ACLs

Weaver ACL is a document formatted in JSON used to decide the destination of downstream traffic. An example of an ACL is
like the following

``` json
{
  "id": "gojek_hello",
  "criterion" : "Method(`POST`) && Path(`/gojek/hello-service`)",
  "endpoint" : {
    "shard_expr": ".serviceType",
    "matcher": "body",
    "shard_func": "lookup",
    "shard_config": {
      "999": {
        "backend_name": "hello_backend",
        "backend":"http://hello.golabs.io"
      }
    }
  }
}
```
The description of each field is as follows

| Field Name |  Description |
|---|---|
| `id`  | The name of the service |
| `criterion`  | The criterion expressed based on [Vulcand Routing](https://godoc.org/github.com/vulcand/route)   |
| `endpoint`  |  The endpoint description (see below) |

For endpoints  the keys descriptions are as following:

| Field Name | Description |
|---|---|
| `matcher` | The value to match can be `body`, `path` , `header` or `param` |
| `shard_expr` | Shard expression, the expression to evaluate request based on the matcher |
| `shard_func` | The function of the sharding (See Below) |
| `shard_config` | The backends for each evaluated value |

For each `shard_config` value there are the value evaluated as the result of expression of `shard_expr`. We need to 
describe backends for each value.

| Field Name | Description |
|---|---|
| `backend_name` | unique name for the evaluated value |
| `backend` | The URI in which the packet will be forwarded |
---
## ACL examples:

Possible  **`shard_func`** values  accepted by weaver are :  `none`, `lookup`, `prefix-lookup`, `modulo` , `hashring`, `s2`. 
Sample ACLs for each  accepted **`shard_func`** are provided below.


**`none`**: 
``` json
{
  "id": "gojek_hello",
  "criterion" : "Method(`GET`) && Path(`/gojek/hello-service`)",
  "endpoint" : {
    "shard_func": "none",
    "shard_config": {
        "backend_name": "hello_backend",
        "backend":"http://hello.golabs.io"
    }
  }
}
```

Details: Just forwards the traffic. HTTP `GET` to `weaver.io/gojek/hello-service` will be forwarded to backend at `http://hello.golabs.io`.

---

**`lookup`**:

``` json
{
  "id": "gojek_hello",
  "criterion" : "Method(`POST`) && Path(`/gojek/hello-service`)",
  "endpoint" : {
    "shard_expr": ".serviceType",
    "matcher": "body",
    "shard_func": "lookup",
    "shard_config": {
      "999": {
        "backend_name": "hello_backend",
        "backend":"http://hello.golabs.io"
      },
      "6969": {
        "backend_name": "bye_backend",
        "backend":"http://bye.golabs.io"
      }
    }
  }
}
```

Details: HTTP `POST` to  `weaver.io/gojek/hello-service` will be forwarded based on the `shard_expr` field's value within request body. The request body shall be similar to: 

``` json 
{
  "serviceType": "999",
  ...
}
```

In this scenario the value evaluated by `shard_expr` will be `999`. This will forward the request to `http://hello.golabs.io`. However if the value evaluated by the `shard_expr` is not found within `shard_config`, weaver returns `503` error. 

---

**`prefix-lookup`**:

``` json 
{
  "id": "gojek_prefix_hello",
  "criterion": "Method(`PUT`) && Path(`/gojek/hello/world`)",
  "endpoint": {
    "shard_expr": ".orderNo",
    "matcher": "body",
    "shard_func": "prefix-lookup",
    "shard_config": {
      "backends": {
          "default": {
              "backend_name": "backend_1",
              "backend": "http://backend1"
          },
          "AB-": {
              "backend_name": "backend2",
              "backend": "http://backend2"
          },
          "AD-": {
              "backend_name": "backend3",
              "backend": "http://backend3"
          }
      },
      "prefix_splitter": "-"
    }
  }
}
```

Details: HTTP `PUT` to `weaver.io/gojek/hello/world` will be forwarded based on the `shard_expr` field's value within request body. The request body shall be similar to: 
``` json 
{
  "orderNo": "AD-2132315",
  ...
}
```
In this scenario the value evaluated by `shard_expr` will be `AD-2132315`. This value will be split according to `prefix_splitter` in `shard_config.backends` evaluating to `AD-`. This will forward the request to `http://backend3`.  However if the value evaluated by the `shard_expr` is not found within `shard_config`, weaver will fallback to the value in the `default` key. If `default` key is not found within `shard_config.backends`, weaver returns `503` error. 

---

**`modulo`**:
``` json
{
  "id": "drivers-by-driver-id",
  "criterion": "Method(`GET`) && PathRegexp(`/v2/drivers/\\d+`)",
  "endpoint": {
    "shard_config": {
      "0": {
        "backend_name": "backend1",
        "backend": "http://backend1"
      },
      "1": {
        "backend_name": "backend2",
        "backend": "http://backend2"
      },
      "2": {
        "backend_name": "backend3",
        "backend": "http://backend3"
      },
      "3": {
        "backend_name": "backend4",
        "backend": "http://backend4"
      }
    },
    "shard_expr": "/v2/drivers/(\\d+)",
    "matcher": "path",
    "shard_func": "modulo"
  }
}
```

Details: HTTP `GET` to `weaver.io/v2/drivers/2156545453242` will be forwarded based on the  value captured by regex in `shard_expr` from `/v2/drivers/2156545453242` path. 
The value must be an integer. The backend is selected based on the modulo operation between extracted value (`2156545453242`) with number of backends in the `shard_config`. In this scenario the result is `2156545453242 % 4 = 2`. This will forward the request to `http://backend3`. 

--- 

**`hashring`**:

``` json
{
  "id": "driver-location",
  "criterion": "Method(`GET`) && PathRegexp(`/gojek/driver/location`)",
  "endpoint": {
    "shard_config": {
      "totalVirtualBackends": 1000,
      "backends": {
        "0-249": {
          "backend_name": "backend1",
          "backend": "http://backend1"
        },
        "250-499": {
          "backend_name": "backend2",
          "backend": "http://backend2"
        },
        "500-749": {
          "backend_name": "backend3",
          "backend": "http://backend3"
        },
        "750-999": {
          "backend_name": "backend4",
          "backend": "http://backend4"
        }
      }
    },
    "shard_expr": "DriverID",
    "matcher": "header",
    "shard_func": "hashring"
  }
}
```

Details: HTTP `PUT` to `weaver.io/gojek/driver/location` will be forwarded based on the result of hashing function from the here.
In this scenario the key by which we select the backend is obtained by using value of DriverID header since matcher is header. For example if request had DriverID: 34345 header, and hashring  calculated hash of that values as hash(34345): 555, it will select backend with 500-749 key. This will forward the request to `http://backend3`

---

**`s2`**:

``` json
{
  "id": "nearby-driver-service-get-nearby",
  "criterion": "Method(`GET`) && PathRegexp(`/gojek/nearby`)",
  "endpoint": {
    "shard_config": {
      "shard_key_separator": ",",
      "shard_key_position": -1,
      "backends": {
        "3477275891585777700": {
          "backend_name": "backend1",
          "backend": "http://backend1",
          "timeout": 300
        },
        "3477284687678800000": {
          "backend_name": "backend2",
          "backend": "http://backend2",
          "timeout": 300
        },
        "3477302279864844300": {
          "backend_name": "backend3",
          "backend": "http://backend3",
          "timeout": 300
        },
        "3477290185236939000": {
          "backend_name": "backend4",
          "backend": "http://backend4",
          "timeout": 300
        }
      }
    },
    "shard_expr": "X-Location",
    "matcher": "header",
    "shard_func": "s2"
  }
}
```

Details: HTTP `GET` to `weaver.io/gojek/nearby` will be forwarded based on the result of s2id calculation from X-Location header in the form of lat and long separated by , in accordance to shard_key_separator value. e.g -6.2428103,106.7940571. Weaver calculated s2id from lat and long because shard_key_position value is -1.

