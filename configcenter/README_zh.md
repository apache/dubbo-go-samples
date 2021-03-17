# 配置中心示例

### 背景

配置中心在 Dubbo 中承担两个职责：

1. 外部化配置。启动配置的集中式存储 （简单理解为 dubbo.properties 的外部化存储）
2. 服务治理。服务治理规则的存储与通知

外部化配置目的之一是实现配置的集中式管理，这部分业界已经有很多成熟的专业配置系统。Dubbo go 中对 Zookeeper、Nacos、Apollo 等做了支持。

外部化配置和其他本地配置在内容和格式上并无区别，配置中心更适合将一些公共配置如注册中心、元数据中心配置等抽取以便做集中管理。

同时，外部化配置有全局和应用两个级别，全局配置是所有应用共享的，应用级配置是由每个应用自己维护且只对自身可见的。并且，外部化配置的优先级默认高于本地配置。

本例以如何在 Zookeeper
上配置和使用外部化配置为例。更多的关于配置中心的描述请参阅 [Dubbo 文档](https://dubbo.apache.org/zh/docs/v2.7/user/configuration/config-center/) 。

### 样例

以 Zookeeper 为例，应用级别和全局的外化配置将会以以下的布局存储在 Zookeeper 上。

```
dubbo
└── config
    ├── dubbo
    │   └── dubbo.properties <- 全局
    ├── user-info-client     
    │   └── dubbo.properties <- 服务消费者
    └── user-info-server     
        └── dubbo.properties <- 服务提供者
```

##### 1. 准备应用外化配置

**全局配置准备** - 全局配置的默认路径为 "/dubbo/config/dubbo/dubbo.properties"。使用以下的命令来配置样例 dubbo.properties：

```bash
zkCli create /dubbo/config; \
zkCli create /dubbo/config/dubbo; \
zkCli create /dubbo/config/dubbo/dubbo.properties; \
zkCli set /dubbo/config/dubbo/dubbo.properties \
"dubbo.protocol.name=dubbo
dubbo.protocol.port=20880"
```

**服务提供方配置** - 假设服务提供方的应用名是 "user-info-server"，服务提供方应用级的配置的默认路径为 "/dubbo/config/user-info-server/dubbo.properties"
。使用以下的命令来配置样例 dubbo.properties：

```bash
zkCli create /dubbo/config; \
zkCli create /dubbo/config/user-info-server; \
zkCli create /dubbo/config/user-info-server/dubbo.properties; \
zkCli set /dubbo/config/user-info-server/dubbo.properties \
"dubbo.service.org.apache.dubbo.UserProvider.cluster=failfast
dubbo.service.org.apache.dubbo.UserProvider.protocol=myDubbo
dubbo.protocols.myDubbo.port=22222
dubbo.protocols.myDubbo.name=dubbo"
```

**服务消费方配置** - 假设服务消费方的应用名是 "user-info-client"，服务提供方应用级的配置的默认路径为 "/dubbo/config/user-info-client/dubbo.properties"
。使用以下的命令来配置样例 dubbo.properties：

```bash
zkCli create /dubbo/config; \
zkCli create /dubbo/config/user-info-client; \
zkCli create /dubbo/config/user-info-client/dubbo.properties; \
zkCli set /dubbo/config/user-info-client/dubbo.properties \
"dubbo.service.org.apache.dubbo.UserProvider.cluster=failfast"
```

##### 2. 在应用中配置配置中心

服务提供方与服务消费方配置配置中心的方式是一致的，如下所示：

```yaml
# config center config
config_center:
  protocol: "zookeeper"
  address: "127.0.0.1:2181"
# application config
application:
  organization: "dubbo.io"
  name: "user-info-server" # 应用文件名，决定应用级外化配置
```

##### 3. 在代码中引入配置中心的包

```go
import (
    _ "github.com/apache/dubbo-go/config_center/zookeeper"
)
```

##### 4. 运行示例

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。


