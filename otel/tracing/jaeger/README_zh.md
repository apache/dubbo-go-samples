# Jaeger 链路追踪导出器

[English](README.md) | [中文](README_zh.md)

本示例展示了 dubbo-go 使用 Jaeger 导出器的链路追踪功能。

## 前置条件

在运行本示例之前，需要先启动一个 Jaeger 实例。可以使用 Docker 运行 Jaeger：

```shell
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

或者使用 docker-compose：

```yaml
version: '3'
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # Jaeger UI
      - "14268:14268"  # Jaeger collector HTTP 端点
    environment:
      - COLLECTOR_OTLP_ENABLED=true
```

## 运行方法

### 启动服务端

```shell
$ go run ./go-server/cmd/main.go
```

### 启动客户端

```shell
$ go run ./go-client/cmd/main.go
```

## 查看追踪数据

运行客户端和服务端后，可以在 Jaeger UI 中查看追踪数据：

1. 在浏览器中打开 `http://localhost:16686`
2. 选择服务名称（例如 `dubbo_otel_jaeger_server` 或 `dubbo_otel_jaeger_client`）
3. 点击 "Find Traces" 查看追踪数据

## 配置说明

Jaeger 导出器通过以下选项进行配置：

- `trace.WithJaegerExporter()`: 使用 Jaeger 作为链路追踪导出器
- `trace.WithEndpoint("http://localhost:14268/api/traces")`: 设置 Jaeger collector 端点
- `trace.WithW3cPropagator()`: 使用 W3C 追踪上下文传播
- `trace.WithAlwaysMode()`: 始终采样追踪数据

如果 Jaeger 实例运行在不同的主机或端口上，可以修改端点地址。

