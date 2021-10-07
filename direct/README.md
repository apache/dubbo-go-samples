# Direct Example

### Backend

In the development and testing environment, it is often necessary to bypass the registry and test the designated service provider, which may require point-to-point direct connection. The point-to-point direct connection method will be based on the service interface and ignore the list of providers in the registry. Interface A is configured point-to-point and does not affect interface B to get the list from the registry.

This example provides the `Consumer` point-to-point direct connection `Provider` based on Dubbo-Go to complete service calls to help better understand the connectivity of Dubbo-Go.

### Introduction

```
├── go-client     
│   ├── cmd       
│   ├── conf      
│   └── pkg         
└── go-server     
    ├── cmd       
    ├── conf      
    ├── docker     
    ├── pkg
    └── tests
        └── integration
```

- go-server: The Service Provider
- go-client: The Service Consumer

#### Provider
Direct example code description:

1. Configure the Dubbo protocol, registry, service information, See [server.yml](go-server/conf/server.yml)

```yaml
services:
  "UserProvider":
    registries: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    methods:
    - name: "GetUser"
      retries: 1
      loadbalance: "random"
```

2. Startup: Register the service

```go
hessian.RegisterPOJO(&pkg.User{})
config.Load()
initSignal()
```

- Based on the `hessian` serialization protocol, using [apache/dubbo-go-hessian2](https://github.com/apache/dubbo-go-hessian2) RegisterPOJO register a POJO
- Dubbo Init: Registration service, See [apache/dubbo-go/../config_loader.go](https://dubbo.apache.org/dubbo-go/v3/blob/master/config/config_loader.go)
    - init router
    - init the global event dispatcher
    - start the metadata report if config set
    - reference config
    - service config
    - init the shutdown callback
- Init Signal ：
  
    ```go
    func initSignal() {
        signals := make(chan os.Signal, 1)
        // It is not possible to block SIGKILL or syscall.SIGSTOP
        signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
        for {
            sig := <-signals
            logger.Infof("get signal %s", sig.String())
            switch sig {
            case syscall.SIGHUP:
                // reload()
            default:
                time.AfterFunc(time.Duration(survivalTimeout), func() {
                    logger.Warnf("app exit now by force...")
                    os.Exit(1)
                })
    
                // The program exits normally or timeout forcibly exits.
                fmt.Println("provider app exit now...")
                return
            }
        }
    }
    ```

#### Consumer

1. Set up the `dubbo service` you need to subscribe to at the beginning of the program startup.
   Make sure that the configuration file [client.yml](go-client/conf/client.yml) has been configured with the relevant information of the subscription service, and the service properties can be customized to override the configuration of the Provider's properties.
   Retain minimum configuration `application` and `references` verification point-to-point direct connection effect, no need to configure the registry.

```go
var userProvider = new(pkg.UserProvider)

func init() {
    config.SetConsumerService(userProvider)
    hessian.RegisterPOJO(&pkg.User{})
}
```

```yaml
application:
  organization: "dubbo.io"
  name: "UserInfoClient"
  module: "dubbo-go user-info client"
  version: "0.0.1"
  environment: "dev"
references:
  "UserProvider":
    registries: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    # this is necessary to enable direct-invoking.
    url: "dubbo://127.0.0.1:20000"
    methods:
      - name: "GetUser"
        retries: 3
```

2. Startup: Direct connection to the service to complete a service call

```go
hessian.RegisterPOJO(&pkg.User{})
config.Load()
user := &pkg.User{}
err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
```


### How To Run

Refer to  [HOWTO.md](../HOWTO_zh.md) under the root directory to run this sample.

#### 1. Environment Configuration

Configure the environment variable to specify the configuration file path required for the service to load.

- go-server:

```shell
APP_LOG_CONF_FILE=direct/go-server/conf/log.yml;
CONF_PROVIDER_FILE_PATH=direct/go-server/conf/server.yml
```

- go-client:

```shell
APP_LOG_CONF_FILE=direct/go-client/conf/log.yml;
CONF_CONSUMER_FILE_PATH=direct/go-client/conf/client.yml
```

See [dubbo-go/.../env.go](https://dubbo.apache.org/dubbo-go/v3/blob/master/common/constant/env.go)


#### 2. Start The Registry

This example uses ZooKeeper as the registry, so you can run the Docker ZooKeeper environment directly. See [docker-compose.yml](go-server/docker/docker-compose.yml)

#### 3. Start The Provider
#### 4. Start The Consumer


Refer to  [HOWTO.md](../HOWTO_zh.md) under the root directory to run this sample.

