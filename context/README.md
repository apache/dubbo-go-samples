# Dubbo-go Context

## 1.Introduction

This sample demonstrates how to use context to pass and read additional parameters.

## 2.How to Pass Additional Parameters
### 2.1 Client

You can pass context parameters to the server by using `context.WithValue` function. Note that the key needs to be fixed to "attachment".

```go
	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.AttachmentKey, map[string]interface{}{
        "key1": "user defined value 1",
        "key2": "user defined value 2"
	})
```

### 2.2 Server

You can get context parameters from the client by using `context.Value` function. Note that the key needs to be fixed to "attachment" and the value type is []string.

```go
    attachments := ctx.Value(constant.AttachmentKey).(map[string]interface{})
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
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/context/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	attachments := ctx.Value(constant.AttachmentKey).(map[string]interface{})
	if value1, ok := attachments["key1"]; ok {
		logger.Infof("Dubbo attachment key1 = %s", value1.([]string)[0])
	}
	if value2, ok := attachments["key2"]; ok {
		logger.Infof("Dubbo attachment key2 = %s", value2.([]string)[0])
	}

	resp := &greet.GreetResponse{Greeting: req.Name}
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

The client file creates a client, writes variables to the context, and makes a call and prints the result.

Source file path：dubbo-go-sample/context/go-client/main.go

```go
package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
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

	ctx := context.Background()
	ctx = context.WithValue(ctx, constant.AttachmentKey, map[string]interface{}{
		"key1": "user defined value 1",
		"key2": "user defined value 2",
	})

	resp, err := svc.Greet(ctx, &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)
}

```

### 3.3 Example Output

Start the server first and then the client. You will observe that the server prints the parameter value passed by the client through the context, indicating that the parameters were successfully passed and obtained.

```
2024-02-26 11:13:14     INFO    logger/logging.go:42    Dubbo attachment key1 = [user defined value 1]
2024-02-26 11:13:14     INFO    logger/logging.go:42    Dubbo attachment key2 = [user defined value 2]
```


