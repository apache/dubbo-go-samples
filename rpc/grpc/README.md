# Dubbo-go GRPC Sample

## 1. Introduction

This sample demonstrates how to achieve interoperability between Dubbo-go and gRPC.

## 2. how to achieve interoperability

Dubbo-go's Triple protocol is compatible with the gRPC protocol. When creating a server, you can set `protocol.WithTriple()` to use the Triple protocol.

```go
srv, err := server.NewServer(
    server.WithServerProtocol(
    protocol.WithPort(20000),
    protocol.WithTriple(),
    ),
)
```

## 3. Example

### 3.1 Server Introduction

#### Server Proto File

Source file path: dubbo-go-sample/rpc/grpc/proto/greet.proto

```protobuf
syntax = "proto3";

package greet;

option go_package = "github.com/apache/dubbo-go-samples/rpc/grpc/proto;greet";

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

#### Dubbo-go Server

Source file path: dubbo-go-sample/rpc/grpc/go-server/cmd/main.go

```go
type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: "dubbo:" + req.Name}
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

#### gRPC Server

Source file path: dubbo-go-sample/rpc/grpc/grpc-server/cmd/main.go

```go
type server struct {
	pb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	resp := &pb.GreetResponse{Greeting: "grpc:" + req.Name}
	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:20001")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})
	logger.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
```

### 3.2 Client Introduction

#### Dubbo-go Client

Source file path: dubbo-go-sample/rpc/grpc/go-client/cmd/main.go

```go
package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/rpc/grpc/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	// Test connection with Dubbo-go
	dubboCli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(dubboCli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)

	// Test connection with gRPC
	grpcCli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20001"),
	)
	if err != nil {
		panic(err)
	}
	svc, err = greet.NewGreetService(grpcCli)
	if err != nil {
		panic(err)
	}
	resp, err = svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)
}
```

#### gRPC Client

Source file path: dubbo-go-sample/rpc/grpc/grpc-client/cmd/main.go

```go
package main

import (
	"context"
	"time"

	"github.com/dubbogo/gost/log/logger"

	pb "github.com/apache/dubbo-go-samples/rpc/grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	// Test connection with gRPC
	conn, err := grpc.Dial("127.0.0.1:20001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreetServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Greet(ctx, &pb.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Fatalf("could not greet: %v", err)
	}
	logger.Infof("Greet response: %s", r.Greeting)
}
```

### Example Output

Start the Dubbo-go server and gRPC server first, then start the Dubbo-go client and gRPC client, and observe the console output.

Dubbo-go client output:

```shell
Greet response: dubbo:hello world
Greet response: grpc:hello world
```

gRPC client output:

```shell
Greet response: grpc:hello world
Greet response: dubbo:hello world
```
