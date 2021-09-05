# dubbogo-grpc

## Contents

- protobuf: proto files for grpc and triple respectively;
- server
- stream-client: clients using bidirectional streaming RPC
- unary-client: clients using unary RPC

Please note that neither server streaming RPC nor client streaming RPC are not supported by Triple so far.

What combination we tested are:

- [x] grpcgo-client(stream) -> dubbogo-server
- [x] grpcgo-client(unary) -> dubbogo-server
- [x] dubbogo-client(stream) -> dubbogo-server
- [x] dubbogo-client(unary) -> dubbogo-server

## Getting Started

### Server

1. Edit your own proto file, please refer to [helloworld.proto](./protobuf/triple/helloworld.proto).
2. Install `protoc` tool, please refer to [ProtoBuf documentation](https://developers.google.com/protocol-buffers/docs/gotutorial).
3. Install `protoc-gen-dubbo3` which is used to generate a stub suitable for triple.

```shell
go get -u github.com/dubbogo/tools/cmd/protoc-gen-triple
```

4. Compile the proto file.

```shell
protoc -I . helloworld.proto --triple_out=plugins=triple:.
```

5. Edit the configuration for server, please refer to [dubbogo.yml](./server/dubbogo-server/conf/dubbogo.yml).
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

2. Edit the configuration for client, please refer to [dubbogo.yml](./stream-client/dubbogo-client/conf/dubbogo.yml)

3. Launch the client.