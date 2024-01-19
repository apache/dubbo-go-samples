# Dubbo-go Timeout Sample

## 1. Introduction

This sample demonstrates how to set a request timeout when creating a Dubbo-go client.

## 2. How to Set Dubbo-go Request Timeout

You can set the timeout for a request by using the `client.WithClientRequestTimeout()` method when creating the client.

```go
    cli, err := client.NewClient(
        client.WithClientURL("tri://127.0.0.1:20000"),
        client.WithClientRequestTimeout(3*time.Second),
    )
```

## 3. Example

### 3.1 Server Introduction

#### Server Proto File

Source file path: dubbo-go-sample/timeout/proto/greet.proto

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

#### Server Handler File

The `Greet` method responds directly, while the `GreetTimeout` method waits for five seconds before responding (simulating a timeout).

Source file path: dubbo-go-sample/timeout/go-server/handler.go

```go
package main

import (
   "context"
   "time"

   greet "github.com/apache/dubbo-go-samples/timeout/proto"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
   resp := &greet.GreetResponse{Greeting: req.Name}
   return resp, nil
}

func (srv *GreetTripleServer) GreetTimeout(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
   time.Sleep(5 * time.Second)
   resp := &greet.GreetResponse{Greeting: req.Name}
   return resp, nil
}
```

### 3.2 Client Introduction

The client file creates a client, sets the timeout to 3 seconds, and makes requests to `Greet` and `GreetTimeout`, outputting the response results.

Source file path: dubbo-go-sample/timeout/go-client/client.go

```go
package main

import (
    "context"
    "time"

    "dubbo.apache.org/dubbo-go/v3/client"
    _ "dubbo.apache.org/dubbo-go/v3/imports"
    greet "github.com/apache/dubbo-go-samples/timeout/proto"
    "github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
		client.WithClientRequestTimeout(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	// test timeout
	resp, err := svc.GreetTimeout(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error("call [greet.GreetService.GreetTimeout] service timeout")
		logger.Error(err)
	} else {
		logger.Infof("Greet response: %s", resp.Greeting)
	}

	// test normal
	resp, err = svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)
}
```

### 3.3 Example Output

Start the server first and then the client. You will observe that the `GreetTimeout` request times out, while the `Greet` request receives a normal response.

```
[call [greet.GreetService.GreetTimeout] service timeout]
Greet response: [hello world]
```