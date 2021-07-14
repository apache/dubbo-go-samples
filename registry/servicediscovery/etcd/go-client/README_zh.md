# 应用级服务发现-etcd示例

### 背景

dubbo-go在提供服务级注册-订阅模式的同时，也提供了应用级服务发现的功能。应用级服务发现使注册中心数据量大大减少。

### 示例

**代码**

在provider和consumer都引入以下package，以初始化etcd注册中心和应用级服务发现功能
```golang
import (
    _ "dubbo.apache.org/dubbo-go/v3/registry/servicediscovery"
    _ "dubbo.apache.org/dubbo-go/v3/registry/etcdv3"
    _ "dubbo.apache.org/dubbo-go/v3/metadata/report/etcd"
)
```

**配置**

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
在本地启动etcd，先后启动provider和consumer，即可验证。
