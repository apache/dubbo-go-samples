# Nacos as registry

This example shows dubbo-go's service discovery feature with Nacos as registry.

## How to run

### Start Nacos server
Follow this instruction to [install and start Nacos server](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

### Run server
```shell
$ go run ./go-server/cmd/server.go
```

test rpc server work as expected:
```shell
$ curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

Open `https://localhost:8848/nacos/` with browser, check url address successfully registered into Nacos.

### Run client
```shell
$ go run ./go-client/cmd/client.go
hello world
```