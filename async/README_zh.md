# 异步 RPC Dubbo for Dubbo-go

[[English](README.md) | [中文](README_zh.md)]

该示例基于新版 `client` / `server` API 展示 Triple 协议下的 Dubbo 异步调用：客户端
在发送 `GetUser` 请求后可以继续执行其他逻辑（非阻塞调用），同时也包含 `SayHello` 的单向调用示例。注意：本示例仅演示异步调用的非阻塞特性，实际响应可通过返回值获取。

## 运行步骤

1. **启动服务端**

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **启动客户端**

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

客户端会打印"non-blocking before async callback resp: do something ... "和"test end"日志， 演示异步调用的非阻塞特性。
