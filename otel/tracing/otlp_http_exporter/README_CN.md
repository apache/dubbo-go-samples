# otlp http tracing exporter

[English](README.md) | 中文

该集成测试检查了 dubbo/triple/jsonrpc 的链路追踪功能，并使用 OpenTelemetry 的 otlp http exporter 来导出追踪数据。

测试会将追踪数据导出到本地的 mock http 服务器。

## 如何运行

### 启动服务端

```shell
go run ./go-server/cmd/main.go
```

### 运行客户端

```shell
go run ./go-client/cmd/main.go
```

如果运行成功，你会在服务端终端看到类似如下的日志：

```shell
2025-09-18 16:22:29	INFO	cmd/main.go:127	server count: 3, client count: 3
```

如果运行失败，你会在服务端终端看到 panic。
