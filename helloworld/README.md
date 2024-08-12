# Helloworld for dubbo-go

This example demonstrates the basic usage of dubbo-go as an RPC framework. Check [Quick Start](https://dubbo.apache.org/zh-cn/overview/mannual/golang-sdk/quickstart/) on our official website for detailed explanation.

## Contents

- go-server/cmd/main.go - is the main definition of the service, handler and rpc server
- go-client/cmd/main.go - is the rpc client
- proto - contains the protobuf definition of the API

## How to run

### Prerequisites
1. Install `protoc` [version3][]
   Please refer to [Protocol Buffer Compiler Installation][].

2. Install `protoc-gen-go` and `protoc-gen-triple`
   Install the version of your choice of protoc-gen-go. here use the latest version as example:

    ```shell
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
    ```
   
    Install the latest version of protoc-gen-triple:

    ```shell
    go install github.com/dubbogo/protoc-gen-go-triple/v3@v3.0.2
    ```

3. Generate stub code

    Generate related stub code with protoc-gen-go and protoc-gen-go-triple:

    ```shell
    protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./proto/greet.proto
    ```


### Run server
```shell
go run ./go-server/cmd/main.go
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
go run ./go-client/cmd/main.go
```

[Quick Start]: https://dubbo-next.staged.apache.org/zh-cn/overview/mannual/golang-sdk/quickstart/
[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler Installation]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/
