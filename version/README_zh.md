# 多版本示例

## 背景

在 Dubbo 中为同一个服务配置多个版本
当一个接口实现，出现不兼容升级时，可以用版本号过渡，版本号不同的服务相互间不引用。


可以按照以下的步骤进行版本迁移：

1. 在低压力时间段，先升级一半提供者为新版本
2. 再将所有消费者升级为新版本
3. 然后将剩下的一半提供者升级为新版本

## 示例介绍

### 目录

```markdown
.
├── README.md
├── README_zh.md
├── docker
├── go-api
├── go-client
├── go-server-v1
└── go-server-v2

```

- go-client ：服务消费者
- go-server-v1 ：服务提供者(老版本)
- go-server-v2 ：服务提供者(新版本)
- go-api ：公共实体

### 介绍

Consume the specified service by changing the Provider and Consumer version numbers

#### 提供者配置

老版本服务提供者配置：

```yaml
# service config
services:
  "UserProvider":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    version: "1.0"
    methods:
    - name: "GetUser"
      retries: 1
      loadbalance: "random"
```

新版本服务提供者配置：

```yaml
# service config
services:
  "UserProvider":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    version: "2.0"
    methods:
    - name: "GetUser"
      retries: 1
      loadbalance: "random"
```

#### 消费者配置

老版本服务消费者配置：

```yaml
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    version: "1.0"
    methods:
      - name: "GetUser"
        retries: 3
```

新版本服务消费者配置：

```yaml
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    version: "2.0"
    methods:
      - name: "GetUser"
        retries: 3
```

#### 效果

消费者日志

```markdown
response result: {A001 Alex Stocks 18 2021-05-09 18:30:16.957 +0800 CST Provider Version 2.0}
```

提供者日志

```markdown
Server V2 req:[]interface {}{"A001"}[2021-05-09/18:30:16 github.com/apache/dubbo-go-samples/version/go-server-v2/pkg.(*UserProvider).GetUser: user.go: 47] %s
Server V2 rsp:pkg.User{ID:"A001", Name:"Alex Stocks", Age:18, Time:time.Time{wall:0xc01e0c4e39144248, ext:20634030436, loc:(*time.Location)(0x4b80960)}, ServiceInfo:"Provider Version 2.0"}
```

### 如何运行

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。

