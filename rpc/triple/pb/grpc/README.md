# Triple <-> gRPC

There are some stuff we provided:

- protobuf: struct definition files for grpc and triple respectively;
- server
- stream-client: clients using streaming RPC
- unary-client: clients using unary RPC

Arbitrary server and client groups can complete RPC calls, as triple is compatible with gRPC. Please note that unidirectional streaming RPC is not supported by triple so far.

## Run Dubbo-go with Triple.

### Server

1. Edit your own proto file, please refer to [helloworld.proto](./protobuf/triple/helloworld.proto).
2. Install `protoc` tool, please refer to [protobuf documentation](https://developers.google.com/protocol-buffers/docs/gotutorial).
3. Install `protoc-gen-dubbo3` which is used to generate a stub suitable for triple.

```shell
    go get -u dubbo.apache.org/dubbo-go/v3/protocol/dubbo3/protoc-gen-dubbo3@3.0
```

4. Compile the proto file.

```shell
    protoc -I . helloworld.proto --dubbo3_out=plugins=grpc+dubbo:.
```

5. Edit the configuration for server, please refer to [dubbogo.yml](./server/dubbogo-server/conf/dubbogo.yml).
6. Launch the server.

### Client

Please note that the start-up process is the same for both the unary RPC and the stream RPC.

1. Register RPC services.

```go
// Directly introduce the GreeterClientImpl structure, you can enter the structure, and see the Reference as "greeterImpl"
var greeterProvider = new(dubbo3pb.GreeterClientImpl)

func init() {
    config.SetConsumerService(greeterProvider)
}
```

2. Edit the configuration for client, please refer to [dubbogo.yml](./stream-client/dubbogo-client/conf/dubbogo.yml)

3. Launch the client.
