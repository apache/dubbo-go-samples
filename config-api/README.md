# config by api

### 1. Usage
provider side\
go-server/cmd/server.go
```go
// run these codes in init function
// and you need not add config env_vairable before run
    providerConfig := config.NewProviderConfig(
            config.WithProviderAppConfig(config.NewDefaultApplicationConfig()),
            config.WithProviderProtocol("dubbo", "dubbo", "20000"),// protocol and port
            config.WithProviderRegistry("demoZk", config.NewDefaultRegistryConfig("zookeeper")), // registry config
            config.WithProviderServices("UserProvider", config.NewServiceConfigByAPI(
                config.WithServiceRegistry("demoZk"), // registry key, equal to upper line
                config.WithServiceProtocol("dubbo"), // export protocol 
                config.WithServiceInterface("org.apache.dubbo.UserProvider"), // interface id
                config.WithServiceLoadBalance("random"), // lb 
                config.WithServiceWarmUpTime("100"),
                config.WithServiceCluster("failover"),
                config.WithServiceMethod("GetUser", "1", "random"),
            )),
        )
	config.SetProviderConfig(*providerConfig) // set to providerConfig ptr
```

consumer side\
go-client/cmd/client.go
```go
consumerConfig := config.NewConsumerConfig(
		config.WithConsumerAppConfig(config.NewDefaultApplicationConfig()), // default app config
		config.WithConsumerConnTimeout(time.Second*3), // timeout
		config.WithConsumerRequestTimeout(time.Second*3), // timeout
		config.WithConsumerRegistryConfig("demoZk", config.NewDefaultRegistryConfig("zookeeper")), // registry config
		config.WithConsumerReferenceConfig("UserProvider", config.NewReferenceConfigByAPI( // set refer config
			config.WithReferenceRegistry("demoZk"), // registry key
			config.WithReferenceProtocol("dubbo"), // protocol 
			config.WithReferenceInterface("org.apache.dubbo.UserProvider"),// interface name
			config.WithReferenceMethod("GetUser", "3", "random"), // method and lb
			config.WithReferenceCluster("failover"),
		)),
	)
	config.SetConsumerConfig(*consumerConfig) // set to global consumerConfig ptr before main function run
```
### 2. Attention
- default registry support\
Now we support the api way to create provider/consumer global configure.

default registry code as showed above:\
we support default ip and port for them
```go
config.NewDefaultRegistryConfig("nacos")
config.NewDefaultRegistryConfig("consul")
config.NewDefaultRegistryConfig("zookeeper")
```

- default app config support\
set app-relevant config by default
```
config.NewDefaultApplicationConfig()
```
- need improve\
Now, api config is only support to create config by api, it also needs to set specific config field after read from yaml file by api.\
There maybe future-change in config/consumer_config.go  config/provider_config.go to export more api.

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.