# Dubbo-go Retry Sample

## 1. Introduction

This sample demonstrates how to use retry feature in Dubbo-go.

## 2. How to use retry feature

When creating a client with `client.NewClient()`, you can use `client.WithClientRetries()` to set the number of retries.

```go
cli, err := client.NewClient(
	client.WithClientURL("tri://127.0.0.1:20000"),
	client.WithClientRetries(3),
)
```

## 3. Example

### 3.1 Server Introduction

#### Server proto file

Source file path: dubbo-go-sample/retry/proto/greet.proto

```protobuf
syntax = "proto3";

package greet;

option go_package = "github.com/apache/dubbo-go-samples/timeout/proto;greet";

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

#### Server handler file

The `Greet` method responds directly, and the `GreetRetry` method is used to simulate retries.

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

### 3.2 Client Introduction

Client file, create client, set retry times to 3, request `Greet` and `GreetRetry` respectively, and observe the server log output.

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

### 3.3 Case effect

First start the server, then start the client, and observe that the server log outputs the number of retries when accessing `GreetRetry`.

```
2024-01-23 22:39:11     INFO    logger/logging.go:22    [Not need retry, request success]
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry [1] times
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry [2] times
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry [3] times
2024-01-23 22:39:11     INFO    logger/logging.go:42    retry success, current request time is [3]
```
