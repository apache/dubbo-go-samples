# Health Check

## Background

Dubbo-go provides a built-in health check service based on the Triple protocol to help users check the health status of their services.

## Usage

- After starting the Dubbo-go framework through an instance, the health check service is automatically registered in the framework, providing health check services based on the Triple protocol without the need for additional configuration in the configuration file.
- The Triple health check service can be used to check the status of services in the framework by making an HTTP request or by invoking the health check service through a client. The service name for the invocation is "grpc.health.v1.Health", and the interface is "check".

## 1. Invoking the Health Check Service through a Client

Start the `instance` service in `dubbo-go-samples/health_check/go-server`. You can use the client code below to check the status of the "greet.GreetService".

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
		client.WithURL("tri://127.0.0.1:20000"),
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

After starting, the following output will be displayed:

```sh
[greet.GreetService's health status:SERVING]
[greet.GreetService's health status:SERVING]
```

## 2. Invoking the Health Check Service by Making an HTTP Request

Start the `instance` service in `dubbo-go-samples/health_check/go-server`. You can make the following HTTP request to check the status of the "greet.GreetService":

```http
POST /grpc.health.v1.Health/Check
Host: 127.0.0.1:20000
Content-Type: application/json

{"service":"greet.GreetService"}
```

The response will be:

```http
{
  "status": "SERVING"
}
```