# Apollo 配置中心示例


## 介绍


### go-server 启动

1. 创建新的apollo服务端yaml格式的namespace

2. 添加服务端配置内容

```yaml
dubbo:
  application:
     name: "demo-server"
     version: "2.0"
  registries:
    "demoZK":
      protocol: "zookeeper"
      timeout: "3s"
      address: "127.0.0.1:2181"
  protocols:
    "triple":
      name: "tri"
      port: 20000
  provider:
    registry-ids:
      - demoZK
    services:
      "greeterImpl":
        protocol-ids: "triple"
        interface: "com.apache.dubbo.sample.basic.IGreeter" # must be compatible with grpc or dubbo-java
```

3. 启动go-server

### go-client 启动

1. 创建新的apollo客户端yaml格式的namespace

2. 添加客户端配置内容

```yaml
dubbo:
  registries:
    "demoZK":
      protocol: "zookeeper"
      timeout: "3s"
      address: "127.0.0.1:2181"
  consumer:
    registry-ids:
      - demoZK
    references:
      "greeterImpl":
        protocol: "tri"
        interface: "com.apache.dubbo.sample.basic.IGreeter" # must be compatible with grpc or dubbo-java
```

3. 启动go-client