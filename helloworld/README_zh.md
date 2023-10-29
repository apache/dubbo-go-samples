# dubbo-go Triple helloworld示例

这个示例旨在帮助您快速地完成一次基础的RPC调用。

## 开始

### 安装版本为[3+](https://protobuf.dev/programming-guides/proto3/)的protoc

详情请参考[**安装Protocol Buffer Compiler**](https://grpc.io/docs/protoc-installation/)和[**下载Protocol Buffers**](https://protobuf.dev/downloads/)。  
安装完成后，请运行```protoc --version``来确保protoc的版本为3+。

### 安装 protoc-gen-go

```shell
# 安装您所需版本的protoc-gen-go，这里使用当前最新版本作为示例
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
```

### 安装 protoc-gen-triple

```shell
# 安装最新版本protoc-gen-triple
git clone https://github.com/apache/dubbo-go.git && cd ./dubbo-go
git checkout feature-triple
go mod tidy
cd ./protocol/triple/triple-tool/protoc-gen-triple
go install .
```

## 生成 Triple stub 代码

```shell
mkdir ~/triple_helloworld && cd ~/triple_helloworld
go mod init triple_helloworld
mkdir proto && cd ./proto

# 编写proto IDL文件
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

# 使用protoc-gen-go和protoc-gen-triple生成相应的stub代码
protoc --go_out=. --go_opt=paths=source_relative --triple_out=. --triple_opt=paths=source_relative ./greet.proto
```

## 完成客户端和服务端代码

### 客户端

编写 **client.go**。具体代码请参考[**这里**](https://github.com/apache/dubbo-go-samples/blob/new-triple-samples/helloworld/client.go)。

### 服务端

实现 **GreetService** 接口并完成 **server.go**。具体代码请参考[**这里**](https://github.com/apache/dubbo-go-samples/blob/new-triple-samples/helloworld/server.go).

## 编译并运行

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