# Triple-Grpc 互通示例

## 示例结构

在本示例中，我们针对 triple 和 grpc ，分别提供了 proto 协议、服务端、普通 RPC 客户端、流式调用而客户端。

可以交叉使用任一组客户端和服务端，实现 RPC 调用。 

## Triple 服务启动

### api 生成
1. 首先编写 proto 文件
  
```protobuf
syntax = "proto3";

option go_package = "protobuf/triple";
package protobuf;

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (User) {}
  rpc SayHelloStream (stream HelloRequest) returns (stream User) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message User {
  string name = 1;
  string id = 2;
  int32 age = 3;
}
```

2. 安装 protoc cli 工具
   
3. 安装 protoc-gen-dubbo3（用于生成适配于triple的stub）
```shell
go get -u dubbo.apache.org/dubbo-go/v3/protocol/dubbo3/protoc-gen-dubbo3@3.0
```
4. 生成 api 文件
```shell
    protoc -I . helloworld.proto --dubbo3_out=plugins=grpc+dubbo:.
```

### 服务端启动
1. Provider 结构定义
```go
package pkg

import (
	"context"
	"fmt"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
)

import (
	dubbo3 "github.com/apache/dubbo-go-samples/general/triple/api/grpc/protobuf/triple"
)

type GreeterProvider struct {
	// 引入dubbo3 provider base
	*dubbo3.GreeterProviderBase
}

func NewGreeterProvider() *GreeterProvider {
	return &GreeterProvider{// 使用生成的pb中的实例初始化base
		GreeterProviderBase: &dubbo3.GreeterProviderBase{},
	}
}

// SayHelloStream 提供流式RPC调用的函数
func (s *GreeterProvider) SayHelloStream(svr dubbo3.Greeter_SayHelloStreamServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 1 user, name = %s\n", c.Name)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 2 user, name = %s\n", c2.Name)
	c3, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 3 user, name = %s\n", c3.Name)

	svr.Send(&dubbo3.User{
		Name: "hello " + c.Name,
		Age:  18,
		Id:   "123456789",
	})
	svr.Send(&dubbo3.User{
		Name: "hello " + c2.Name,
		Age:  19,
		Id:   "123456789",
	})
	return nil
}

// SayHello 提供普通rpc调用的服务函数
func (s *GreeterProvider) SayHello(ctx context.Context, in *dubbo3.HelloRequest) (*dubbo3.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n" + in.Name)
	fmt.Println("get triple header tri-req-id = ", ctx.Value("tri-req-id"))
	fmt.Println("get triple header tri-service-version = ", ctx.Value("tri-service-version"))
	return &dubbo3.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}

// Reference 需要和config中Reference key 对应
func (g *GreeterProvider) Reference() string {
	return "GreeterProvider"
}
```
2. 配置文件
server/dubbogo-server/conf/server.yml
```yaml
# service config
services:
  "GreeterProvider":
    registry: "demoZK"
    protocol: "tri" # 使用triple协议
    interface: "protobuf.Greeter" # 和grpc生成的的接口名一致，如下
```

Grpc api 文件的接口名可见为 protobuf.Greeter, 是用户根据需要定义的。
triple-go 要想和 grpc 打通，一定要和 grpc 的接口名一致，并正确配置在 yaml 文件中。

protobuf/grpc/helloworld.api.go:
```go
func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/protobuf.Greeter/SayHello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
```

3. 启动服务端

goland 运行
triple/triple-server

## 客户端（以普通rpc调用为例， 流式RPC同理）

1. Consumer 端结构定义
Consumer 结构已经在pb文件中实现好,可直接引入
```go
import (
    dubbo3pb "github.com/apache/dubbo-go-samples/general/triple/api/grpc/protobuf/triple"
)

// 直接引入GreeterClientImpl 结构，可以进入该结构，查看Reference为“greeterImpl”
var greeterProvider = new(dubbo3pb.GreeterClientImpl)

func init() {
    config.SetConsumerService(greeterProvider)
}
```

配置文件中定义好对应 reference key
```yaml
# reference config
references:
  "greeterImpl":
    registry: "demoZk"
    protocol: "tri"
    interface: "protobuf.Greeter"
    url: tri://localhost:20001
```

2. 启动客户端

goland 运行
triple/triple-unary-client

