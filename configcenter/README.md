# Config Center Sample

### Background

In Dubbo, config center takes the following two responsibilities:
1. For externalized configuration. Central management for bootstrap configurations, that is, keep "dubbo.properties" save externally on the config center.
2. For service governance. Keep service rule on the config center, and notify the subscribers when the rule changes.

The purpose of externalized configuration is centrally managing application's configurations in one single place. There are many mature configuration management systems in the industry. Dubbo go supports some of them, including Zookeeper, Nacos and Apollo, etc.

There's no difference between externalized configuration and local configuration with regarding the content and the format. Config center is more suitable to keep the common configurations such as registry center configuration, metadata center configuration, etc., so that they are centrally managed.

At the same time, externalized configuration has two scope: global scope, and application scope. Global scope configurations are shared among all applications, and application scope configuration is only visible to the application it belongs to. By default, externalized configuration has higher priority than the local configuration.

This sample demonstrates how to config and use externalized configuration on Zookeeper. You can refer [Dubbo Documentation](https://dubbo.apache.org/zh/docs/v2.7/user/configuration/config-center/) for more details on config center.

### Example

If use Zookeeper as config center, global externalized configuration and application scope configurations will be layout like below:

```
dubbo
└── config
    ├── dubbo
    │   └── dubbo.properties <- global
    ├── user-info-client     
    │   └── dubbo.properties <- service consumer
    └── user-info-server     
        └── dubbo.properties <- service provider
```

##### 1. Prepare externalized configurations

**Prepare global scope externalized configuration** - The default path for the global scope configuration on Zookeeper is: "/dubbo/config/dubbo/dubbo.properties". You can use the following commands to config it:

```bash
zkCli create /dubbo/config; \
zkCli create /dubbo/config/dubbo; \
zkCli create /dubbo/config/dubbo/dubbo.properties; \
zkCli set /dubbo/config/dubbo/dubbo.properties \
"dubbo.protocol.name=dubbo
dubbo.protocol.port=20880"
```

**Prepare configuration for service provider** - Assume service provider's application name is "user-info-server", then the default path for the provider's config is "/dubbo/config/user-info-server/dubbo.properties". You can use the following commands to config it:

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

**Prepare configuration for service consumer** - Assume service consumer's application name is "user-info-client", then the default path for the consumer's config is "/dubbo/config/user-info-client/dubbo.properties". You can use the following commands to config it:

```bash
zkCli create /dubbo/config; \
zkCli create /dubbo/config/user-info-client; \
zkCli create /dubbo/config/user-info-client/dubbo.properties; \
zkCli set /dubbo/config/user-info-client/dubbo.properties \
"dubbo.service.org.apache.dubbo.UserProvider.cluster=failfast"
```

##### 2. Config config center in application

The config center's configuration on the consumer side and on the provider side is similar, shown as following:

```yaml
# config center config
config_center:
  protocol: "zookeeper"
  address: "127.0.0.1:2181"
# application config
application:
  organization: "dubbo.io"
  name: "user-info-server" # application's name, which will decide where the externalized configuration is placed.
```

##### 3. Import config center's package

```go
import (
    _ "github.com/apache/dubbo-go/config_center/zookeeper"
)
```

##### 4. Run this sample

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.


