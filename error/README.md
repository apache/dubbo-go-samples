# Dubbo-go Error

## 1.Introduction

This sample demonstrates how to handle calling error in dubbo-go 

## 2.Example

### 2.1 Server

#### Server Proto File

Source file path：dubbo-go-sample/error/proto/greet.proto

```protobuf
syntax = "proto3";

package greet;

option go_package = "github.com/apache/dubbo-go-samples/error/proto;greet";

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

In this program, only when the request parameter "name" accepted by the Greet function is "right name", it is considered a correct request, otherwise it is considered an error request and will not be processed and an error will be returned.

Source file path：dubbo-go-sample/context/go-server/main.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	greet "github.com/apache/dubbo-go-samples/helloworld/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	name := req.Name
	if name != "right name" {
		errInfo := fmt.Sprintf("name is not right: %s", name)
		logger.Error(errInfo)
		return nil, errors.New(errInfo)
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

	if err = greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}

```

### 2.2 Client

The client initiates two calls respectively, one is a correct request and one is an incorrect request.

Source file path：dubbo-go-sample/context/go-client/main.go

```go
package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/helloworld/proto"
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

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "right name"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("call Greet success: %s", resp.Greeting)

	resp, err = svc.Greet(context.Background(), &greet.GreetRequest{Name: "wrong name"})
	if err != nil {
		logger.Errorf("call Greet failed, err: %s", err.Error())
	}
}
```

### 2.3 Example Output

Start the server first, and then the client. We can observe that the first call is successful, but the second call fails and the server returns an error.

```
2024-02-28 17:49:40	INFO	logger/logging.go:42	call Greet success: [right name]
2024-02-28 17:49:40	ERROR	logger/logging.go:52	call Greet failed, err: [Failed to invoke the method Greet in the service greet.GreetService. Tried 2 times of the providers [tri://:@127.0.0.1:20000/?interface=greet.GreetService&group=&version= tri://:@127.0.0.1:20000/?interface=greet.GreetService&group=&version= tri://:@127.0.0.1:20000/?interface=greet.GreetService&group=&version=] (3/1)from the registry tri://127.0.0.1:20000/greet.GreetService?app.version=&application=dubbo.io&async=false&bean.name=greet.GreetService&cluster=failover&config.tracing=&environment=&generic=&group=&interface=greet.GreetService&loadbalance=&metadata-type=local&module=sample&name=dubbo.io&organization=dubbo-go&owner=dubbo-go&peer=true&provided-by=&reference.filter=cshutdown&registry.role=0&release=dubbo-golang-3.2.0&remote.timestamp=&retries=&serialization=protobuf&side=consumer&sticky=false&timestamp=1709113780&version= on the consumer 30.221.146.234 using the dubbo version 3.2.0. Last error is unknown: name is not right: wrong name.: unknown: name is not right: wrong name]
```


