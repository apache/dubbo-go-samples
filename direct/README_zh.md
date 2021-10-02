# Direct 示例

### 背景

在开发及测试环境下，经常需要绕过注册中心，测试指定服务提供者，这时候可能需要点对点直连，
点对点直连方式，将以服务接口为单位，忽略注册中心的提供者列表，A 接口配置点对点，不影响 B 接口从注册中心获取列表。

本示例提供基于 Dubbo-go 的`Consumer`点对点直连`Provider`完成服务调用，帮助更好理解 Dubbo-go 的连通性。


### 目录

```
├── go-client     
│   ├── cmd       启动入口
│   ├── conf      消费者配置：dubbo 服务属性配置、日志属性配置
│   └── pkg       业务包  
└── go-server     
    ├── cmd       启动入口
    ├── conf      服务提供者配置：dubbo 服务属性配置、日志属性配置
    ├── docker    docker compose： Zookeeper 
    ├── pkg
    └── tests
        └── integration
```

- go-server: 服务提供者
- go-client: 服务消费者

#### 服务提供者

直连示例代码说明：

1. 配置dubbo 服务协议，注册中心，服务信息等，具体参阅 [server.yml](go-server/conf/server.yml)

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

2. 应用启动：注册服务

```go
hessian.RegisterPOJO(&pkg.User{})
config.Load()
initSignal()
```

- 基于`hessian`序列化协议，使用[apache/dubbo-go-hessian2](https://github.com/apache/dubbo-go-hessian2) RegisterPOJO注册一个POJO实例
- Dubbo Init：注册服务，详情参考[apache/dubbo-go/../config_loader.go](https://dubbo.apache.org/dubbo-go/v3/blob/master/config/config_loader.go)
    - init router
    - init the global event dispatcher
    - start the metadata report if config set
    - reference config
    - service config
    - init the shutdown callback
- 初始化 signal 包：将输入信号（对应信号）转发到 `chan`， signal包不会为了向`chan`发送信息而阻塞，调用者应该保证`chan`有足够的缓存空间可以跟上期望的信号频率，此处单一信号用于通知的通道，缓存为设置 `1`
    
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

#### 服务消费者

1. 在程序启动之初设置需要订阅的 `dubbo` 服务，
   确保配置文件 `client.yml` 已配置订阅服务相关信息，可自定义设置服务属性等，覆盖 Provider 的属性配置，详情参阅 [client.yml](go-client/conf/client.yml),
   保留最少配置 `application` 和 `references` 验证点对点直连效果，无需注册中心等配置

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

2. 应用启动：直连服务，完成一次服务调用

```go
hessian.RegisterPOJO(&pkg.User{})
config.Load()
user := &pkg.User{}
err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
```


### 如何运行
请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。
#### 1. 环境配置

配置环境变量，指定服务加载所需配置文件路径

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

详情请参阅 [dubbo-go/.../env.go](https://dubbo.apache.org/dubbo-go/v3/blob/master/common/constant/env.go)


#### 2. 启动注册中心

本示例使用 Zookeeper 做注册中心， 可以直接运行 docker zookeeper 环境，配置详情请参阅 `docker-compose.yml`

#### 3. 启动服务提供者
#### 4. 启动服务消费者


请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。

