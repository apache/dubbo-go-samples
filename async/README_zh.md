# 异步 RPC Dubbo for Dubbo-go

[[English](README.md) | [中文](README_zh.md)]

该示例基于新版 `client` / `server` API 展示 Triple 协议下的 Dubbo 异步调用：客户端
在发送 `GetUser` 请求后可以继续执行其他逻辑，随后通过回调拿到服务端的结果，同时也
包含 `SayHello` 的单向调用示例。

## 运行步骤

1. **启动服务端**

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **启动客户端**

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

客户端会先打印“非阻塞”日志，随后在收到回调结果时再次打印响应内容。
