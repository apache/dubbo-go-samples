# 游戏服务示例

### 背景介绍

- 示例包含 **gate** (网关服务) 和 **game** (逻辑服务) 两个服务
- 两个服务会互相 RPC 通讯 (都同时注册 **provider** 和 **consumer**)
- **gate** 额外启动了 http 服务 (端口 **8000**), 用于手工触发  **gate** RPC 调用 **game**

> 每次 **gate** RPC调用(**Message**) **game** 后, **game** 会同步RPC调用(Send) **gate** 推送相同消息

### 概要

目录说明

```bash
├── go-server-game    # game模块
│   ├── cmd           # 主入口
│   ├── conf          # 配置文件
│   ├── docker        # docker-compose文件
│   ├── pkg           # provider和consumer
│   └── tests
├── go-server-gate    # gate模块
│   ├── cmd           # 主入口
│   ├── conf          # 配置文件
│   ├── docker        # docker-compose文件
│   ├── pkg           # provider和consumer
│   └── tests
└── pkg
    ├── consumer      # 公共consumer
    │   ├── game
    │   └── gate
    └── pojo
```

发起http服务的流程

<img src="http://cdn.cjpa.top/cdnimages/image-20210423212453935.png" alt="image-20210423212453935" style="zoom:50%;" />



从consumer和provider角度来看，发起一次调用的流程是这样的

<img src="http://cdn.cjpa.top/cdnimages/image-20210424094134541.png" alt="image-20210424094134541" style="zoom: 33%;" />

game提供了basketball服务端，gate提供了http服务端。

### game模块

#### server端

server 端提供三个服务，Login、Score 及 Rank，代码如下，具体的实现可以在 'game/go-server-game/pkg/provider.go' 中看到

```go
type BasketballService struct{}

func Login(ctx context.Context, data string) (*pojo.Result, error) {
    ...
}

func Score(ctx context.Context, uid, score string) (*pojo.Result, error) {
    ...
}

func Rank (ctx context.Context, uid string) (*pojo.Result, error) {
    ...
}

func (p *BasketballService) Reference() string {
    return "gameProvider.basketballService"
}
```

##### 配置文件

 在配置文件中注册 service，其中 gameProvider.basketballService 和 Reference 方法中声明的一致。

```yml
services:
  "gameProvider.basketballService":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.game.BasketballService"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    methods:
      - name: "Online"
        retries: 0
      - name: "Offline"
        retries: 0
      - name: "Message"
        retries: 
```

#### consumer端

basketball 部分的 consumer 主要用来被 gate 调用，因此主要代码放在了 'game/pkg/consumer/gate/basketball.go' 中，这部分作为公共部分，consumer 代码如下

```go
type BasketballService struct {
    Send func(ctx context.Context, uid string, data string) (*pojo.Result, error)
}

func (p *BasketballService) Reference() string {
    return "gateConsumer.basketballService"
}
```

在basketball中，只需要实例化一个consumer变量即可

```go
var gateBasketball = new(gate.BasketballService)
```

然后在main中注册到dubbo-go

```go
config.SetConsumerService(gateBasketball)
```

##### 配置文件

```yml
references:
  "gateConsumer.basketballService":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.gate.BasketballService"
    cluster: "failover"
    methods:
      - name: "Send"
        retries: 0
```

由于 game 的 consumer 需要调用 gate 的provider，因此 Reference 方法返回的字符串为 gateConsumer.basketballService ，在配置文件中 gateConsumer.basketballService 和 'game/pkg/consumer/gate/basketball.go' 中Reference 方法声明的一致，intreface 的值也要和 gate 的 provider 设置的一致。

### gate模块

#### server端

```go
type BasketballService struct{}

func (p *BasketballService) Send(ctx context.Context, uid, data string) (*pojo.Result, error) {
...
}

func (p *BasketballService) Reference() string {
    return "gateProvider.basketballService"
}
```

注册到dubbo

```go
config.SetProviderService(new(BasketballService))
```

##### 配置文件

gateProvider.basketballService 和 Reference 中的一致，这里的 interface 一定要和 game 的 client.yml 文件中设置的保持一致，不然 game 无法向 gate 发送数据

```yml
# service config
services:
  "gateProvider.basketballService":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.gate.BasketballService"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    methods:
      - name: "Send"
        retries: 0
```

#### consumer端

gate 中的 consumer 端比较特殊，由于 gate的consumer 需要调用 game 中的 service，所以在gaet中，consumer 直接实例化一个game的service，其方法便直接使用实例化的对象 GameBasketball 调用，这样就实现了一个网关的功能。

```go
var GameBasketball = new(game.BasketballService)

func Login(ctx context.Context, data string) (*pojo.Result, error) {
    return GameBasketball.Login(ctx, data)
}

func Score(ctx context.Context, uid, score string) (*pojo.Result, error) {
    return GameBasketball.Score(ctx, uid, score)
}

func Rank (ctx context.Context, uid string) (*pojo.Result, error) {
    return GameBasketball.Rank(ctx, uid)
}
```

代码中的 GameBasketball.Message、GameBasketball.Online、GameBasketball.Offline 调用的方法都是 game 的方法

注册到dubbo

```go
config.SetProviderService(new(pkg.BasketballService))
```

##### 配置文件

配置文件中的 inerface 也要和 game 的 provider 保持一致，不然收到 http 请求之后无法调用 gate

```yml
references:
  "gameConsumer.basketballService":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.game.BasketballService"
    cluster: "failover"
    methods:
      - name: "Online"
        retries: 0
      - name: "Offline"
        retries: 0
      - name: "Message"
        retries: 0
```



### HTTP访问

访问login

```bash
curl --location --request GET 'http://127.0.0.1:8089/login?name=dubbogo'
```

收到账户的信息

```json
{
    "code": 0,
    "msg": "dubbogo, your score is 0",
    "data": {
        "score": 0,
        "to": "dubbogo"
    }
}
```



访问 score

```bash
curl --location --request POST 'http://127.0.0.1:8089/score' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"cjp",
    "score":1
}'
```

收到

```json
{
    "code": 0,
    "msg": "dubbogo, your score is 0",
    "data": {
        "score": 0,
        "to": "dubbogo"
    }
}
```

访问 rank

```http
curl --location --request POST 'http://127.0.0.1:8089/rank?name=dubbogo'
```

收到

```json
{
    "code": 0,
    "msg": "success",
    "data": {
        "rank": 3,
        "to": "dubbogo"
    }
}
```

可以发现，所有的http请求，都成功的转发到了 game，并且又通过调用 gate 的 send 方法把处理过的数据返回给 gate



### 启动前端

切换到 webside 目录，打开 index.html 即可（推荐使用chrome浏览器）

![image-20210516173728198](http://cdn.cjpa.top/image-20210516173728198.png)

点击小人即可进行游戏

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。

