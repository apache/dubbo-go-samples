# dubbo-go Triple helloworld示例

这个示例旨在帮助你快速地完成一次基础的RPC调用。

## 开始

### 安装protoc

```shell
# 根据您的操作系统和架构下载预先编译好的protoc文件
# protoc-<version>-<os>-<arch>.zip
# 或者您可以从 github.com/protocolbuffers/protobuf/releases 手动下载
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v24.4/protoc-24.4-linux-x86_64.zip

# 在一个目录下解压下载好的protoc文件，这里使用$HOME/.local作为示例
unzip protoc-24.4-linux-x86_64.zip -d $HOME/.local

# 更新环境变量
export PATH="$PATH:$HOME/.local/bin"
```

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

```shell
cd ~/triple_helloworld
mkdir -p go-client/cmd && cd ./go-client/cmd
```

编写 **client.go** 并把它放置在 **go-client/cmd** 目录。

```go
package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/client"
	// 重要，使用这个import声明引用dubbo-go扩展
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "triple_helloworld/proto"
	"triple_helloworld/proto/greettriple"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	// 初始化一个用于调用特定服务的Client，如果你想要调用其它服务，请创建新的Client
	cli, err := client.NewClient(
		// 指定该服务的URL
		client.WithURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}
	
	svc, err := greettriple.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp.Greeting)
}
```

### 服务端
```shell
cd ~/triple_helloworld
mkdir -p go-server/cmd
mkdir -p go-server/handler
```

实现 **GreetService** 接口并把它放置在 **go-server/handler** 目录。
请参考[**具体实现**](https://github.com/apache/dubbo-go-samples/tree/new-triple-samples/helloworld/go-server/handler).

完成 **server.go** 并把它放置在 **go-server/cmd** 目录。

```go
package main

import (
	// 重要，使用这个import声明引用dubbo-go扩展
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"triple_helloworld/go-server/handler"
	"triple_helloworld/proto/greettriple"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	// 初始化一个Server用于承载多个服务
	srv, err := server.NewServer(
		// 默认使用Triple协议
		server.WithServerProtocol(
			// 指定监听的端口
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	// 注册一个特定服务
	if err := greettriple.RegisterGreetServiceHandler(srv, &handler.GreetTripleServer{}); err != nil {
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
```

## 编译并运行

```shell
cd ~/triple_helloworld/go-server/cmd
go build -o server .
./server
```

```shell
cd ~/triple_helloworld/go-client/cmd
go build -o client .
./client
```