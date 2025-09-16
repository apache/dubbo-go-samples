# 健康检查

## 背景

Dubbo-go 内置了基于Triple协议的健康检查服务，帮助用户检测服务健康状态。

## 使用方法

- Dubbo-go框架在通过instance启动后会自动向框架中注册健康检查服务，提供基于triple服务的健康检查服务，无需在配置文件中额外配置。
- triple健康检查服务可以通过发起http请求检查框架中服务的状态，也可以通过客户端调用该健康检查服务，调用的接口为“grpc.health.v1.Health”，方法为“check”。

## 1、通过客户端调用健康检查服务

启动dubbo-go-samples/health_check/go-server中的服务，通过下方客户端即可查看“greet.GreetService”的状态。

```go
package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	health "dubbo.apache.org/dubbo-go/v3/protocol/triple/health/triple_health"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}
	svc, err := health.NewHealth(cli)
	if err != nil {
		panic(err)
	}
	check, err := svc.Check(context.Background(), &health.HealthCheckRequest{Service: "greet.GreetService"})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("greet.GreetService's health", check.String())
	}
	watch, err := svc.Watch(context.Background(), &health.HealthCheckRequest{Service: "greet.GreetService"})
	if err != nil {
		logger.Error(err)
	} else {
		if watch.Recv() {
			logger.Info("greet.GreetService's health", watch.Msg().String())
		}
	}
}
```

启动后会有以下输出

```sh
[greet.GreetService's health status:SERVING]
[greet.GreetService's health status:SERVING]
```

## 2.通过发起http请求调用健康检查服务

启动dubbo-go-samples/health_check/go-server中的服务，发起下方http请求即可查看“greet.GreetService”的状态

```http
POST /grpc.health.v1.Health/Check
Host: 127.0.0.1:20000
Content-Type: application/json

{"service":"greet.GreetService"}
```

将会有以下输出

```http
{
  "status": "SERVING"
}
```

