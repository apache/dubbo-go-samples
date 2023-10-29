# Helloworld for dubbo-go Triple

This is **Triple** helloworld example to help you finish a basic RPC invocation done quickly.

## Prerequisites

### install protocol buffer compiler, protoc, [version3](https://protobuf.dev/programming-guides/proto3/)

Please refer to [**Protocol Buffer Compiler Installation**](https://grpc.io/docs/protoc-installation/) and [**Download Protocol Buffers**](https://protobuf.dev/downloads/).  
After installation, please run ```protoc --verion``` to ensure that the version of protoc is 3+.

### install protoc-gen-go

```shell
# install the version of your choice of protoc-gen-go. here use the latest version as example
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
```

### install protoc-gen-triple

```shell
# install the latest version of protoc-gen-triple
git clone https://github.com/apache/dubbo-go.git && cd ./dubbo-go
git checkout feature-triple
go mod tidy
cd ./protocol/triple/triple-tool/protoc-gen-triple
go install .
```

## Generate Triple stub code

```shell
mkdir ~/triple_helloworld && cd ~/triple_helloworld
go mod init triple_helloworld
mkdir proto && cd ./proto

# replace this with your own proto IDL file
cat > greet.proto << EOF
syntax = "proto3";

package greet;

option go_package = "triple_helloworld/proto;greet";

message GreetRequest {
  string name = 1;
}

message GreetResponse {
  string greeting = 1;
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
}

EOF

# generate related stub code with protoc-gen-go and protoc-gen-triple
protoc --go_out=. --go_opt=paths=source_relative --triple_out=. --triple_opt=paths=source_relative ./greet.proto
```

## Finish client and server code

### client

Finish **client.go**. For specific code, please refer to [**this**](https://github.com/apache/dubbo-go-samples/blob/new-triple-samples/helloworld/client.go)

### server

Implement **GreetService** Interface and finish **server.go**. For specific code, please refer to [**this**](https://github.com/apache/dubbo-go-samples/blob/new-triple-samples/helloworld/server.go).

## Build and run

```shell
cd ~/triple_helloworld
go build -o server ./server.go
./server
```

```shell
cd ~/triple_helloworld
go build -o client ./client.go
./client
```