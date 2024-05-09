# Dubbo-go Retry Sample

## 1.介绍

本示例演示如何在Dubbo-go中使用重试功能。

## 2.如何使用重试功能

在使用`client.NewClient()`创建客户端时，可以使用`client.WithClientRetries()`方法设置重试次数。

```go
cli, err := client.NewClient(
	client.WithClientURL("tri://127.0.0.1:20000"),
	client.WithClientRetries(3),
)
```

## 3.案例

### 3.1服务端介绍

#### 服务端proto文件

源文件路径：dubbo-go-sample/retry/proto/greet.proto

```protobuf
syntax = "proto3";

package greet;

option go_package = "github.com/apache/dubbo-go-samples/retry/proto;greet";

message GreetRequest {
  string name = 1;
}

message GreetResponse {
  string greeting = 1;
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
  rpc GreetTimeout(GreetRequest) returns (GreetResponse) {}
}
```

#### 服务端handler文件

`Greet`方法直接响应，`GreetRetry`方法用于模拟重试。
```go
package main

import (
	"context"
	"github.com/pkg/errors"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/retry/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
	requestTime int
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	logger.Info("Not need retry, request success")
	return resp, nil
}

func (srv *GreetTripleServer) GreetRetry(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	if srv.requestTime < 3 {
		srv.requestTime++
		logger.Infof("retry %d times", srv.requestTime)
		return nil, errors.New("retry")
	}
	resp := &greet.GreetResponse{Greeting: req.Name}
	logger.Infof("retry success, current request time is %d", srv.requestTime)
	srv.requestTime = 0
	return resp, nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{
		requestTime: 0,
	}); err != nil {
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}

```

### 3.2客户端介绍

客户端client文件，创建客户端，设置重试次数为3，分别请求`Greet`和`GreetRetry`，观察服务端日志输出。

```go
package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/retry/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
		client.WithClientRetries(3),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	// request normal
	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)

	// request need retry
	resp, err = svc.GreetRetry(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)
}
```

### 3.3案例效果

先启动服务端，再启动客户端，访问`GreetRetry`时观察到服务端日志输出了重试次数。

```
2024-01-23 22:39:11     INFO    logger/logging.go:22    [Not need retry, request success]
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry [1] times
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry [2] times
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry [3] times
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry success, current request time is [3]
```