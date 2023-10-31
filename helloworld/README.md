# Helloworld for dubbo-go

This example demonstrates the basic usage of dubbo-go as an RPC framework. Check [Quick Start][] on our official website for detailed explanation.

## Contents

- server/main.go - is the main definition of the service, handler and rpc server
- client/main.go - is the rpc client
- proto - contains the protobuf definition of the API

## How to run

[//]: # (### Prerequisites)

[//]: # ()
[//]: # (1. Install `protoc [version3][]`)

[//]: # ()
[//]: # (    Please refer to [Protocol Buffer Compiler Installation][].)

[//]: # ()
[//]: # (2. Install `protoc-gen-go` and `protoc-gen-triple` )

[//]: # ()
[//]: # (    ```shell)

[//]: # (    # install the version of your choice of protoc-gen-go. here use the latest version as example)

[//]: # (    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31)

[//]: # (    )
[//]: # (    # install the latest version of protoc-gen-triple)

[//]: # (    git clone https://github.com/apache/dubbo-go.git && cd ./dubbo-go)

[//]: # (    git checkout feature-triple)

[//]: # (    go mod tidy)

[//]: # (    cd ./protocol/triple/triple-tool/protoc-gen-triple)

[//]: # (    go install .)

[//]: # (    ```)

[//]: # ()
[//]: # (### Generate stub code)

[//]: # ()
[//]: # (```shell)

[//]: # (# generate related stub code with protoc-gen-go and protoc-gen-triple)

[//]: # (protoc --go_out=. --go_opt=paths=source_relative --triple_out=. --triple_opt=paths=source_relative ./greet.proto)

[//]: # (```)

### Run server
```shell
go run ./server/main.go
```

test server work as expected:
```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

### Run client
```shell
go run ./client/main.go
```

[Quick Start]: https://dubbo-next.staged.apache.org/zh-cn/overview/mannual/golang-sdk/quickstart/
[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler Installation]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/