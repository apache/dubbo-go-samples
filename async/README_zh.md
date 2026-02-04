# Dubbo-go 异步 RPC

[[English](README.md) | [中文](README_zh.md)]

本示例展示了如何使用新的 `client`/`server` API 通过 Triple 协议异步调用 Dubbo 服务。
演示了 Go 和 Java 之间的异步互操作。

## 功能特性

- **Go 客户端和服务端**: 使用 `client.WithAsync()` 实现异步调用
- **Java 客户端**: 使用 `CompletableFuture` API 实现异步调用
- **Java 服务端**: 使用 `CompletableFuture` 实现异步服务
- **互操作性**: Java 客户端可调用 Go 服务端，Go 客户端可调用 Java 服务端

## 运行 Go 到 Go 示例

1. **启动 Go 服务端**

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **启动 Go 客户端**（默认连接 Go 服务端）

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

客户端会打印 "non-blocking before async callback resp: do something ... " 和 "test end" 日志，演示异步调用的非阻塞特性。

## 运行 Java-Go 互操作示例

演示**跨语言异步调用**：

- **Go 客户端** → **Java 服务端**
- **Java 客户端** → **Go 服务端**

### 前置条件

- Java 11 或更高版本
- Maven 3.6+

### 构建 Java 模块

在 `async` 目录下执行：

```bash
mvn clean compile
```

### 测试：Go 客户端 → Java 服务端

1. **修改 Go 客户端 URL**。Go 客户端默认配置为连接 Go 服务端（端口 `20000`），在 `go-client/cmd/main.go` 中通常如下：

   ```go
   client.WithClientURL("tri://127.0.0.1:20000"),
   ```

   为了让 Go 客户端连接到 Java 服务端，请将上述行修改为：

   ```go
   client.WithClientURL("tri://127.0.0.1:50051"),
   ```

2. **启动 Java 服务端**（端口 50051）

   ```bash
   cd java-server
   ./run.sh
   ```

3. **启动 Go 客户端**

   ```bash
   go run ./async/go-client/cmd/main.go
   ```

Go 客户端会向 Java 服务端发送异步请求，并打印 "non-blocking before async callback resp: do something ... " 日志。

### 测试：Java 客户端 → Go 服务端

1. **启动 Go 服务端**（端口 20000）

   ```bash
   go run ./async/go-server/cmd/main.go
   ```

2. **启动 Java 客户端**

   ```bash
   cd java-client
   ./run.sh
   ```

Java 客户端会向 Go 服务端发送异步请求，使用 `CompletableFuture` 回调处理响应。

## 端口分配

- **Go 服务端**: 20000
- **Java 服务端**: 50051

两个服务端可以同时运行，不会产生端口冲突。
