# dubbogo-grpc

## Contents

- protobuf: 使用 proto 文件的结构体定义
- server
- stream-client: 使用流式调用的客户端
- unary-client: 使用普通调用的客户端

请注意，到目前为止，Triple还不支持服务器流式RPC和客户端流式RPC。

我们测试的组合包括:

- [x] grpcgo-client(stream) -> dubbogo-server
- [x] grpcgo-client(unary) -> dubbogo-server
- [x] dubbogo-client(stream) -> dubbogo-server
- [x] dubbogo-client(unary) -> dubbogo-server

## 运行

### 服务端

1. 编辑你自己的 proto 文件，请参考 [helloworld.proto](./protobuf/triple/helloworld.proto)。
2. 安装 `protoc` 工具，请参考 [ProtoBuf 文档](https://developers.google.com/protocol-buffers/docs/gotutorial)。
3. 安装 `protoc-gen-dubbo3`，用于生成适用于 triple 的 stub。

```shell
go get -u github.com/dubbogo/tools/cmd/protoc-gen-triple
```

4. 编译 proto 文件。

```shell
protoc -I . helloworld.proto --triple_out=plugins=triple:.
```

5. 编辑服务端配置文件，请参考 [dubbogo.yml](go-server/conf/dubbogo.yml)。
6. 启动服务端。

### 客户端

请注意，普通调用的 RPC 和流式 RPC 的启动过程是相同的。

1. 注册 RPC 服务。

```go
// Directly introduce the GreeterClientImpl structure, you can enter the structure, and see the Reference as "greeterImpl"
var greeterProvider = new(triplepb.GreeterClientImpl)
func init() {
    config.SetConsumerService(greeterProvider)
}
```

2. 编辑客户端配置文件，请参考 [dubbogo.yml](go-client/conf/dubbogo.yml)。

3. 启动客户端。