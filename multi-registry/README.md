# Multi Registry Example

## Backend

dubbo-go supports simultaneous registration of the same service with multiple registries, or separate registration of different services in different registries, or even simultaneous reference of services with the same name registered in different registries. In addition, the registry supports custom extensions. I use `zk` and `nacos` in this sample.

### Code

**Config**

```yaml
# registry config
registries:
  "zk":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2181"
  "nacos":
    protocol: "nacos"
    timeout	: "3s"
    address: "127.0.0.1:8848"

# reference config
references:
  "UserProvider":
    registry: "zk,nacos"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    methods:
      - name: "GetUser"
        retries: 3
```

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.
