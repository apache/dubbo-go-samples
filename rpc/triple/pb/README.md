# dubbogo-triple-pb

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
protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. greet.proto
```
