# Application-level service discovery - etcd sample

### background

When providing register-subscribe mode for user, dubbo-go also provides application-level service discovery, which can greatly reduce data amount on registry.
### Sample

**Code**

You should import packages below to initiate etcd service discovery.

```golang
import (
    _ "dubbo.apache.org/dubbo-go/v3/registry/servicediscovery"
    _ "dubbo.apache.org/dubbo-go/v3/registry/etcdv3"
    _ "dubbo.apache.org/dubbo-go/v3/metadata/report/etcd"
)
```

**Config**

server.yml

```yaml
# application config
application:
  #               #
  # other config  #
  #               #

  metadataType: "remote"

# registry config
registries:
  "demoServiceDiscovery":
    protocol: "service-discovery"
    params:
      service_discovery: "etcdv3"
      name_mapping: "dynamic"
      metadata: "default"

metadata_report:
  protocol: "etcdv3"
  remote_ref: "etcd"

service_discovery:
  etcdv3:
    protocol: "etcdv3"
    remote_ref: "etcd"
    group: "dubbo"

remote:
  etcd:
    address: "127.0.0.1:2379"
    timeout: "5s"
```

client.yml

```yaml
# application config
application:
  #               #
  # other config  #
  #               #

  metadataType: "remote"

# registry config
registries:
  "demoServiceDiscovery":
    protocol: "service-discovery"
    params:
      service_discovery: "etcdv3"
      name_mapping: "in-memory"
      metadata: "default"

metadata_report:
  protocol: "etcdv3"
  remote_ref: "etcd"

service_discovery:
  etcdv3:
    protocol: "etcdv3"
    remote_ref: "etcd"
    group: "dubbo"

remote:
  etcd:
    address: "127.0.0.1:2379"
    timeout: "5s"
```
Start etcd locally with port 2379, then start provider and consumer.Try it out!
