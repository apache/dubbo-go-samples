# Multi Registry 多注册中心

## 背景

dubbo-go 支持同一服务向多注册中心同时注册，或者不同服务分别注册到不同的注册中心上去，甚至可以同时引用注册在不同注册中心上的同名服务。另外，注册中心是支持自定义扩展的。下面我使用zk和nacos演示

### 代码

**配置**

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


请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。
