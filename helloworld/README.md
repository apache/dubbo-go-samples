# Helloworld for dubbo-go Triple

This is **Triple** helloworld example to help you finish a basic RPC invocation done quickly.

## Prerequisites

### install protoc

```shell
# fetch the pre-compiled protoc corresponding to your operating system and computer architecture
# protoc-<version>-<os>-<arch>.zip
# or you can download from github.com/protocolbuffers/protobuf/releases manually
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v24.4/protoc-24.4-linux-x86_64.zip

# unzip downloaded file under a directory
unzip protoc-24.4-linux-x86_64.zip -d $HOME/.local

# update path variable
export PATH="$PATH:$HOME/.local/bin"
```

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

```shell
cd ~/triple_helloworld
mkdir -p go-client/cmd && cd ./go-client/cmd
```

Finish **client.go** and put it in **go-client/cmd** directory.

```go
package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/client"
	// important, must import this for dubbo-go extensions
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "triple_helloworld/proto"
	"triple_helloworld/proto/greettriple"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	// initialize a Client which is responsible for invoking a certain service. it uses Triple protocol by default
	// if you want to invoke another service, please initialize a new one
	cli, err := client.NewClient(
		// specify target server URL
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

### server
```shell
cd ~/triple_helloworld
mkdir -p go-server/cmd
mkdir -p go-server/handler
```

Implement **GreetService** Interface and put it in **go-server/handler** directory.
Please refer to [**concrete sample**](https://github.com/apache/dubbo-go-samples/tree/new-triple-samples/helloworld/go-server/handler).

Finish **server.go** and put it in **go-server/cmd** directory.

```go
package main

import (
    // important, must import this for dubbo-go extensions
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"triple_helloworld/go-server/handler"
	"triple_helloworld/proto/greettriple"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	// initialize a Server for serving multiple services 
	srv, err := server.NewServer(
		// use Triple protocol by default
		server.WithServerProtocol(
			// specify port to listen on
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	// register a certain service
	if err := greettriple.RegisterGreetServiceHandler(srv, &handler.GreetTripleServer{}); err != nil {
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
```

## Build and run

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