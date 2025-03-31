# Dubbo-go Context

## 1.Introduction

This sample demonstrates how to use context to pass and read additional parameters.

## 2.How to Pass Additional Parameters
### 2.1 Client

You can pass context parameters to the server by using `triple_protocol.AppendToOutgoingContext` function. Note that the key needs to be fixed to "attachment".

```go
    header := http.Header{"testKey1": []string{"testVal1"}, "testKey2": []string{"testVal2"}}
    // to store outgoing data ,and reserve the location for the receive field.
    // header will be copy , and header's key will change to be lowwer.
	ctx := triple_protocol.NewOutgoingContext(context.Background(), header)
    ctx = triple_protocol.AppendToOutgoingContext(ctx, "testKey3", "testVal3")
```

### 2.2 Server

You can get context parameters from the client by using `triple.FromIncomingContext` function. Note that the key is the key confirmed by the sender and the value type is []string.

```go
    data, _ := triple.FromIncomingContext(ctx)
    logger.Infof("Dubbo attachment key1 = %s", value1.([]string)[0])
    logger.Infof("Dubbo attachment key2 = %s", value2.([]string)[0])
```

## 3.Example

### 3.1 Server Introduction

#### Server Proto File

Source file path：dubbo-go-sample/context/proto/greet.proto

```protobuf
syntax = "proto3";

package greet;

option go_package = "github.com/apache/dubbo-go-samples/context/proto;greet";

message GreetRequest {
  string name = 1;
}

message GreetResponse {
  string greeting = 1;
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}
```

#### Server Handler File

Source file path：dubbo-go-sample/context/go-server/main.go

```go
package main

import (
	"context"
	triple "dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"fmt"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/context/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	data, _ := triple.FromIncomingContext(ctx)
	ctx = triple.AppendToOutgoingContext(ctx, "OutgoingContextKey1", "OutgoingDataVal1", "OutgoingContextKey2", "OutgoingDataVal2")
	var value1, value2, value3 string
	if values, ok := data["testkey1"]; ok && len(values) > 0 {
		value1 = values[0]
		logger.Infof("testkey1: %s", value1)
	}
	if values, ok := data["testkey2"]; ok && len(values) > 0 {
		value2 = values[0]
		logger.Infof("testkey2: %s", value2)
	}
	if values, ok := data["testkey3"]; ok && len(values) > 0 {
		value3 = values[0]
		logger.Infof("testkey3: %s", value3)
	}

	respStr := fmt.Sprintf("name: %s, testKey1: %s, testKey2: %s", req.Name, value1, value2)
	resp := &greet.GreetResponse{Greeting: respStr}
	return resp, nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
```

### 3.2 Client Introduction

The client file creates a client, writes variables to the ctx created by `triple_protocol.NewOutgoingContext`, and makes a call and prints the result.

Source file path：dubbo-go-sample/context/go-client/main.go

```go
package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/protocol/triple/triple_protocol"
	"net/http"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/context/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	header := http.Header{"testKey1": []string{"testVal1"}, "testKey2": []string{"testVal2"}}
	// to store outgoing data ,and reserve the location for the receive field.
	// header will be copy , and header's key will change to be lowwer.
	ctx := triple_protocol.NewOutgoingContext(context.Background(), header)
	ctx = triple_protocol.AppendToOutgoingContext(ctx, "testKey3", "testVal3")

	resp, err := svc.Greet(ctx, &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	extractedHeader, _ := triple_protocol.FromIncomingContext(ctx)

	var value1, value2 string
	if values, ok := extractedHeader["outgoingcontextkey1"]; ok && len(values) > 0 {
		value1 = values[0]
		logger.Infof("OutgoingContextKey1: %s", value1)
	}
	if values, ok := extractedHeader["outgoingcontextkey2"]; ok && len(values) > 0 {
		value2 = values[0]
		logger.Infof("OutgoingContextKey2: %s", value2)
	}

	logger.Infof("Greet response: %s", resp.Greeting)
}
```

### 3.3 Example Output

Start the server first and then the client. You will observe that the server prints the parameter value passed by the client through the context, indicating that the parameters were successfully passed and obtained.

server result:
```
2025-03-31 16:55:00     INFO    logger/logging.go:42    testkey1: testVal1
2025-03-31 16:55:00     INFO    logger/logging.go:42    testkey2: testVal2
2025-03-31 16:55:00     INFO    logger/logging.go:42    testkey3: testVal3
```
client result:
```
2025-03-31 17:10:09     INFO    logger/logging.go:42    OutgoingContextKey1: OutgoingDataVal1
2025-03-31 17:10:09     INFO    logger/logging.go:42    OutgoingContextKey2: OutgoingDataVal2
2025-03-31 17:10:09     INFO    logger/logging.go:42    Greet response: name: hello world, testKey1: testVal1, testKey2: testVal2
```

