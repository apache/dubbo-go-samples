# Game Service Example

### Backend

- The example includes two services, **gate** (gateway service) and **game** (logical service)
- The two services communicate with each other RPC (both registered **provider** and **consumer**)
- The **gate** additionally starts the http service (port **8000**), which is used to manually trigger the **gate** RPC to call **game**

> Each time **gate** RPC calls (**message**） **game**, **game** will synchronously call RPC (send) **gate** to push the same message. 

### outline

Catalog description

```shell
├── go-server-game    # game module
│   ├── cmd           # Main entrance
│   ├── conf          # configuration file
│   ├── docker        # docker-compose config file
│   ├── pkg           # provider and consumer
│   └── tests
├── go-server-gate    # gate module
│   ├── cmd           # Main entrance
│   ├── conf          # configuration file
│   ├── docker        # docker-compose config file
│   ├── pkg           # provider and consumer
│   └── tests
└── pkg
    ├── consumer      # public consumer
    │   ├── game
    │   └── gate
    └── pojo
```

Process of initiating HTTP service

<img src="http://cdn.cjpa.top/cdnimages/image-20210424095907886.png" alt="image-20210424095907886" style="zoom: 33%;" />



From the perspective of consumer and provider, the process of initiating a call is as follows.

<img src="http://cdn.cjpa.top/cdnimages/image-20210424100148028.png" alt="image-20210424100148028" style="zoom:33%;" />

Game provides the basketball server and gate provides the HTTP server.

### game module

#### server side

The server provides three services, message, online and offline. The code is as follows. The specific implementation can be in  game/go-server-game/pkg/provider.go see in.

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

##### configuration file

Register the service in the configuration file, where gameProvider.basketballService Same as declared in the reference method.

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

#### consumer side

The consumer of basketball is mainly used to be called by gate, so the main code is put in 'game/pkg/consumer/gate/basketball.go',  The consumer code is as follows.

```go
type BasketballService struct {
    Send func(ctx context.Context, uid string, data string) (*pojo.Result, error)
}

func (p *BasketballService) Reference() string {
    return "gateConsumer.basketballService"
}
```

In basketball, you only need to instantiate a consumer variable.

```go
var gateBasketball = new(gate.BasketballService)
```

Then register to Dubbo go in main.

```go
config.SetConsumerService(gateBasketball)
```

##### configuration file

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

Because the consumer of game needs to call the provider of gate, the string returned by the reference method is gateConsumer.basketballService , in the configuration file gateConsumer.basketballService And 'game/pkg/consumer/gate/basketball.go'. The value of interface should be consistent with that of the provider of gate.

### gate module

#### server side

```go
type BasketballService struct{}

func (p *BasketballService) Send(ctx context.Context, uid, data string) (*pojo.Result, error) {
...
}

func (p *BasketballService) Reference() string {
    return "gateProvider.basketballService"
}
```

Register with Dubbo

```go
config.SetProviderService(new(BasketballService))
```

##### configuration file

gateProvider.basketballService  The interface here must be the same as that in game client.yml  The settings in the file should be consistent, otherwise game cannot send data to gate.

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

#### consumer side

The consumer side of the gate is relatively special. Because the consumer of the gate needs to call the service in the game, the consumer directly instantiates the service of a game, and its method directly uses the instantiated object gamebasketball to call, thus realizing the function of a gateway.

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

In code GameBasketball.Message、GameBasketball.Online、GameBasketball.Offline , all the methods called are game methods.
Register with dubbo

```go
config.SetProviderService(new(pkg.BasketballService))
```

##### configuration file

The inerface in the configuration file should also be consistent with the game provider, otherwise the gate cannot be called after receiving the HTTP request.

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

### HTTP access

access login

```bash
curl --location --request GET 'http://127.0.0.1:8089/login?name=dubbogo'
```

received

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



access score

```bash
curl --location --request POST 'http://127.0.0.1:8089/score' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"cjp",
    "score":1
}'
```

received

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

access rank

```bash
curl --location --request POST 'http://127.0.0.1:8089/rank?name=dubbogo'
```

received

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

It can be found that all HTTP requests are successfully forwarded to game, and the processed data is returned to gate by calling the send method of gate



### Start the front end

Switch to webside directory and open index.html (Chrome browser is recommended)

![image-20210516173728198](http://cdn.cjpa.top/image-20210516173728198.png)

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.

