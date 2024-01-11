# dubbogo-grpc

## Contents

- protobuf: proto files for grpc and triple respectively;
- go-server: Dubbo-go server
- go-client: Dubbo-go client
- grpc-server: gRPC server
- grpc-client: gRPC client

Please note that neither server streaming RPC nor client streaming RPC are not supported by Triple so far.

What combinations we tested are:

- [x] grpcgo-client -> dubbogo-server
- [x] dubbogo-client -> dubbogo-server

## Getting Started

### Server

1. Edit your own proto file, please refer to [samples_api.proto](/api/samples_api.proto).
2. Install `protoc` tool, please refer to [ProtoBuf documentation](https://developers.google.com/protocol-buffers/docs/gotutorial).
3. Install `protoc-gen-dubbo3` which is used to generate a stub suitable for triple.

```shell
go get -u github.com/dubbogo/tools/cmd/protoc-gen-triple
```

4. Compile the proto file.

```shell
protoc -I . helloworld.proto --triple_out=plugins=triple:.
```

5. Edit the configuration for server, please refer to [dubbogo.yml](go-server/conf/dubbogo.yml).
6. Launch the server.

### Client

Please note that the start-up process is the same for both the unary RPC and the stream RPC.

1. Register RPC services.

```go
// Directly introduce the GreeterClientImpl structure, you can enter the structure, and see the Reference as "greeterImpl"
var greeterProvider = new(triplepb.GreeterClientImpl)
func init() {
    config.SetConsumerService(greeterProvider)
}
```

2. Edit the configuration for client, please refer to [dubbogo.yml](go-client/conf/dubbogo.yml)

3. Launch the client.