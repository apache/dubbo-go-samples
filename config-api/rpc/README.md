# Dubbogo 使用api进行配置初始化

### 使用方法
- 服务端：
  triple/go-server/cmd/server.go

```go
// 注册Provider，通过结构名，或者 Reference()函数返回值'greeterImpl'，与 provider配置serviceKey相对应
config.SetProviderService(&GreeterProvider{})

// 生成service配置
serviceConfig := config.NewServiceConfig(
    config.WithServiceInterface("com.apache.dubbo.sample.basic.IGreeter"), // 接口名，consuemr与provider对应的ID
    config.WithServiceProtocolKeys("tripleKey"), // 选择的协议key，与rootConfig 声明的支持协议key需要对应
)

// 生成协议配置
protocolConfig := config.NewProtocolConfig(
    config.WithProtocolName("tri"),
    config.WithProtocolPort("20000"),
)

// 生成provider配置
providerConfig := config.NewProviderConfig(
    config.WithProviderRegistryKeys("zk"), // provider 选择使用的注册中心key
    config.WithProviderService("greeterImpl", serviceConfig), // provider 注册 serviceKey 和用上面定义好的 GreeterProvider
)

// 生成注册中心配置
registryConfig := config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")

// 生成根配置
rootConfig := config.NewRootConfig(
    config.WithRootProviderConfig(providerConfig), // 使用provider配置
    config.WithRootRegistryConfig("zk", registryConfig), // 声明注册中心Key和注册中心配置
    config.WithRootProtocolConfig("tripleKey", protocolConfig),// 声明当前应用支持的协议key，与 serviceConfig 选择的协议key需要对应
)

// 服务启动
if err := rootConfig.Init(); err != nil {
    panic(err)
}
select {}
```
- 客户端，与provider端类似
```go
config.SetConsumerService(tripleGreeterImpl)

referenceConfig := config.NewReferenceConfig(
  config.WithReferenceInterface("com.apache.dubbo.sample.basic.IGreeter"),
  config.WithReferenceProtocolName("tri"),
  config.WithReferenceRegistry("zkRegistryKey"),
)

consumerConfig := config.NewConsumerConfig(
  config.WithConsumerReferenceConfig("greeterImpl", referenceConfig),
)

registryConfig := config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")

rootConfig := config.NewRootConfig(
  config.WithRootRegistryConfig("zkRegistryKey", registryConfig),
  config.WithRootConsumerConfig(consumerConfig),
)

if err := rootConfig.Init(); err != nil {
  panic(err) 
}
```