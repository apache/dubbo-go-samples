# Dubbo-go Context

## 1.介绍

本示例演示如何在Dubbo-go框架中如何传递和读取附加参数，来实现上下文信息传递。

## 2.使用说明
### 2.1客户端使用说明

在客户端中，使用下述 `triple_protocol.AppendToOutgoingContext` 方式传递字段:

```go
    header := http.Header{"testKey1": []string{"testVal1"}, "testKey2": []string{"testVal2"}}
    // to store outgoing data ,and reserve the location for the receive field.
    // header will be copy , and header's key will change to be lowwer.
	ctx := triple_protocol.NewOutgoingContext(context.Background(), header)
    ctx = triple_protocol.AppendToOutgoingContext(ctx, "testKey3", "testVal3")
```

### 2.2服务端使用说明

在服务端中，使用下述方式获取`triple.FromIncomingContext`字段, key为发送方传入的key, value的类型为[]string:
```go
    data, _ := triple.FromIncomingContext(ctx)
    logger.Infof("Dubbo attachment key1 = %s", value1.([]string)[0])
    logger.Infof("Dubbo attachment key2 = %s", value2.([]string)[0])
```

## 3.案例

### 3.1服务端介绍

#### 服务端proto文件

源文件路径：dubbo-go-sample/context/proto/greet.proto

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

#### 服务端handler文件

源文件路径：dubbo-go-sample/context/go-server/main.go

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

### 3.2客户端介绍

客户端client文件，创建客户端，在`triple_protocol.NewOutgoingContext`所创建的ctx里写入变量，发起调用并打印结果

源文件路径：dubbo-go-sample/context/go-client/main.go

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

### 3.3案例效果

先启动服务端，再启动客户端，可以观察到服务端打印了客户端传递的参数值，说明参数被成功传递并获取

服务端打印结果:
```
2025-03-31 16:55:00     INFO    logger/logging.go:42    testkey1: testVal1
2025-03-31 16:55:00     INFO    logger/logging.go:42    testkey2: testVal2
2025-03-31 16:55:00     INFO    logger/logging.go:42    testkey3: testVal3
```
客户端打印结果:
```
2025-03-31 17:10:09     INFO    logger/logging.go:42    OutgoingContextKey1: OutgoingDataVal1
2025-03-31 17:10:09     INFO    logger/logging.go:42    OutgoingContextKey2: OutgoingDataVal2
2025-03-31 17:10:09     INFO    logger/logging.go:42    Greet response: name: hello world, testKey1: testVal1, testKey2: testVal2
```


