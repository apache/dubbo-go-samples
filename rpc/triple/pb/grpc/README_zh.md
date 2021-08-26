# Triple <-> gRPC

我们提供的案例包括：

- protobuf: 使用 proto 文件的结构体定义
- server
- stream-client: 使用流式调用的客户端
- unary-client: 使用普通调用的客户端

任意的服务器和客户端组可以完成RPC调用，因为 Triple 与 gRPC 兼容。请注意，到目前为止 Triple 还不支持单向流 RPC。

## 运行

### 服务端

1. 编辑你自己的 proto 文件，请参考 [helloworld.proto](./protobuf/triple/helloworld.proto)。
2. 安装 `protoc` 工具，请参考 [protobuf documentation](https://developers.google.com/protocol-buffers/docs/gotutorial)。
3. 安装 `protoc-gen-dubbo3`，用于生成适用于 triple 的 stub。

```shell
    go get -u dubbo.apache.org/dubbo-go/v3/protocol/dubbo3/protoc-gen-dubbo3@3.0
```

4. 编译 proto 文件。

```shell
    protoc -I . helloworld.proto --dubbo3_out=plugins=grpc+dubbo:.
```

5. 编辑服务端配置文件，请参考 [dubbogo.yml](./server/dubbogo-server/conf/dubbogo.yml)。
6. 启动服务端。

### 客户端

请注意，普通调用的 RPC 和流式 RPC 的启动过程是相同的。

1. 注册 RPC 服务。

```go
// Directly introduce the GreeterClientImpl structure, you can enter the structure, and see the Reference as "greeterImpl"
var greeterProvider = new(dubbo3pb.GreeterClientImpl)

func init() {
    config.SetConsumerService(greeterProvider)
}
```

2. 编辑客户端配置文件，请参考 [dubbogo.yml](./stream-client/dubbogo-client/conf/dubbogo.yml)。

3. 启动客户端。