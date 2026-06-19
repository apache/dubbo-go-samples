# Dubbo-Go Triple Header 和 Trailer 示例

中文 | [English](README.md)

本示例演示如何在 Triple 协议中传递请求元数据，以及读写响应 header / trailer。

示例覆盖：

- 客户端通过 `triple.NewOutgoingContext` 和 `triple.AppendToOutgoingContext` 发送请求 header
- 服务端通过 `triple.FromIncomingContext` 和生成流式接口的 `RequestHeader` 读取请求 metadata
- Unary 调用通过 `triple.SetHeader`、`triple.SetTrailer`、`client.WithResponseHeader` 和 `client.WithResponseTrailer` 读写响应 header / trailer
- 服务端通过流式接口的 `ResponseHeader` / `ResponseTrailer` 设置响应 header / trailer
- 客户端通过生成的流式接口读取响应 header / trailer
- 按 `http.Header` 的方式使用 `Values("X-Sample-Token")` 读取 header

这和 `context` 示例不同。`context` 示例演示的是通过 `constant.AttachmentKey` 传递 Dubbo attachments；本示例演示的是 Triple metadata API，并以 `http.Header` 的形式暴露给应用代码。

对于 Dubbo-Go 生成的 Unary 调用，服务端通过 `triple.SetHeader` 和 `triple.SetTrailer` 写入响应 metadata，客户端通过 `client.WithResponseHeader` 和 `client.WithResponseTrailer` 捕获本次调用的响应 metadata。对于生成的双向流和客户端流调用，应在创建 stream 前通过 `NewOutgoingContext` / `AppendToOutgoingContext` 传入请求 metadata。服务端 handler 可以通过 `triple.FromIncomingContext` 读取这些 metadata；生成的 stream handler 也可以通过 `RequestHeader` 查看。

## 运行方式

启动服务端：

```bash
go run ./triple_header_trailer/go-server/cmd
```

在另一个终端启动客户端：

```bash
go run ./triple_header_trailer/go-client/cmd
```

也可以通过集成测试脚本运行：

```bash
./integrate_test.sh triple_header_trailer
```

## 重新生成代码

`proto` 目录下的生成代码来自 `proto/greet.proto`：

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./triple_header_trailer/proto/greet.proto
```

## 说明

Triple metadata key 语义上大小写不敏感。在 Go API 层，这些 metadata 以 `http.Header` 暴露，因此应用代码应按标准 `http.Header` 方式使用 `Get` / `Values` 读取。
