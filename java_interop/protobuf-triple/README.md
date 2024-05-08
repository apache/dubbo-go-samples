# Dubbo java and go interoperability, protobuf and triple protocol

使用同一个 proto 文件实现 dubbo-java 和 dubbo-go 互通

## Contents

- go，go 语言实现的 rpc server 与 client
- java，java 语言实现的 rpc server 与 client

## 互通模式

共享服务定义如下：

```protobuf
//protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. greet.proto
syntax = "proto3";
package org.apache.dubbo.sample;

option go_package = "github.com/apache/dubbo-go-samples/java_interop/protobuf-triple/go/proto;proto";
//package of go
option java_package = 'org.apache.dubbo.sample';
option java_multiple_files = true;
option java_outer_classname = "HelloWorldProto";
option objc_class_prefix = "WH";

// The greeting service definition.
service Greeter {
  // Sends a greeting
  rpc SayHello(HelloRequest) returns (HelloReply);
  // Sends a greeting via stream
  //  rpc SayHelloStream (stream HelloRequest) returns (stream HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

### Java client 调用 go server

1. 首先启动 go server：

```shell
go run go/go-server/cmd/server.go
```

运行以上命令后，go server 运行在 50052 端口，可通过以下命令测试服务运行正常：

```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:50052/org.apache.dubbo.sample.Greeter/sayHello
```

2. 启动 java client 

运行以下命令，启动 java 客户端，可以看到服务调用 go server 正常输出结果：

```shell
#在java-client目录中
./run.sh
```

### Go client 调用 java server

1. 启动 java server

运行以下命令，启动 java 服务端：

> 注意，请关闭之前启动的 go server，避免出现端口占用冲突。

```shell
#在java-server目录中
./run.sh
```

可通过以下命令测试服务运行正常：

```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:50052/org.apache.dubbo.sample.Greeter/sayHello
```

2. 运行 go client

运行以下命令启动 go 客户端，可以看到服务调用 java server 正常输出结果：

```shell
go run go/go-client/cmd/client.go
```
