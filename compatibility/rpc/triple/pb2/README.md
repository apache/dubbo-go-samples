# dubbogo-grpc(pb2)

This sample is a simple example of dubbo-go-grpc(pb2) with triple protocol.
It's using `go-to-protobuf` to generate pb2 files from Go struct.


## Contents

- api: proto files for grpc and triple respectively;
- go-server: Dubbo-go server
- go-client: Dubbo-go client
- models: models for Go server and client
- hack: hack scripts for generating

Please note that neither server streaming RPC nor client streaming RPC are not supported by Triple so far.


## Build and Run

1. Install dev-tools

```shell
go install k8s.io/code-generator/cmd/go-to-protobuf@latest
go install github.com/gogo/protobuf/protoc-gen-gogo@latest
go install github.com/dubbogo/tools/cmd/protoc-gen-go-triple@v1.0.8
go install github.com/golang/protobuf/protoc-gen-go@latest
```

2. Generate pb files and go files

```shell
# NOTICE: make sure project in your $GOPATH/src , cause go-to-protobuf will use $GOPATH/src as proto file path
# using vendor as proto path
go mod vendor

# generate pb files from go struct
bash rpc/triple/pb2/hack/gen-go-to-protobuf.sh

# generate RPC go files from pb files
protoc \
  --proto_path=. \
  --proto_path="$GOPATH/src" \
  --go_out=rpc/triple/pb2/api \
  --go-triple_out=rpc/triple/pb2/api \
  rpc/triple/pb2/api/helloworld.proto
  
# remove vendor
rm -rf vendor
```

3. Run

```shell
# start a zk as a registry
docker run --rm --name some-zookeeper -p 2181:2181 zookeeper

# start server
DUBBO_GO_CONFIG_PATH=$(pwd)/rpc/triple/pb2/go-server/conf/dubbogo.yml go run rpc/triple/pb2/go-server/cmd/server.go

# start client
DUBBO_GO_CONFIG_PATH=$(pwd)/rpc/triple/pb2/go-client/conf/dubbogo.yml go run rpc/triple/pb2/go-client/cmd/client.go
```