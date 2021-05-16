# Dubbogo 使用api进行配置初始化

### 1. 使用方法
- 服务端：
  go-server/cmd/server.go

```go
// 在init函数中定义这些代码
// 你不需要在环境变量中定义服务端配置文件的位置了
    providerConfig := config.NewProviderConfig(
            config.WithProviderAppConfig(config.NewDefaultApplicationConfig()), // 默认app配置
            config.WithProviderProtocol("dubbo", "dubbo", "20000"),// 协议key、协议名和端口号
            config.WithProviderRegistry("demoZk", config.NewDefaultRegistryConfig("zookeeper")), // 注册中心配置
            config.WithProviderServices("UserProvider", config.NewServiceConfigByAPI(
                config.WithServiceRegistry("demoZk"), // 注册中心 key, 和上面注册中心key保持一致
                config.WithServiceProtocol("dubbo"), // 暴露协议，和上面协议key对应
                config.WithServiceInterface("org.apache.dubbo.UserProvider"), // interface id
                config.WithServiceLoadBalance("random"), // 负载均衡
                config.WithServiceWarmUpTime("100"),
                config.WithServiceCluster("failover"),
                config.WithServiceMethod("GetUser", "1", "random"),
            )),
        )
	config.SetProviderConfig(*providerConfig) // 写入providerConfig指针
```

- 客户端
  go-client/cmd/client.go

```go
consumerConfig := config.NewConsumerConfig(
		config.WithConsumerAppConfig(config.NewDefaultApplicationConfig()), // 默认app配置
		config.WithConsumerConnTimeout(time.Second*3), // timeout
		config.WithConsumerRequestTimeout(time.Second*3), // timeout
		config.WithConsumerRegistryConfig("demoZk", config.NewDefaultRegistryConfig("zookeeper")), // 注册中心配置
		config.WithConsumerReferenceConfig("UserProvider", config.NewReferenceConfigByAPI( // set refer config
			config.WithReferenceRegistry("demoZk"), // registry key
			config.WithReferenceProtocol("dubbo"), // protocol 
			config.WithReferenceInterface("org.apache.dubbo.UserProvider"),// interface name
			config.WithReferenceMethod("GetUser", "3", "random"), // method and lb
			config.WithReferenceCluster("failover"),
		)),
	)
	config.SetConsumerConfig(*consumerConfig) // 写入 consumerConfig 指针
```
### 2. 注意
- 默认注册中心支持

  现在提供根据api来快速初始化注册中心到默认ip和端口

默认注册中心初始化代码:

```go
config.NewDefaultRegistryConfig("nacos")
config.NewDefaultRegistryConfig("consul")
config.NewDefaultRegistryConfig("zookeeper")
```

- 默认app配置
设置默认的app配置：
```
config.NewDefaultApplicationConfig()
```
- 需要改进的地方：

  现在，通过api只能进行配置初始化过程，而通过api进行配置特定字段的修改也很重要。在以后可能会在config/consumer_config.go  config/provider_config.go 文件中进行进一步支持。

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。