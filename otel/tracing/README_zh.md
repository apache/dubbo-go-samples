# 链路追踪

[English](./README.md) | 中文

## 背景

Dubbo-go 支持 [OpenTelemetry](https://opentelemetry.io/) 链路追踪。

- [Stdout 导出器](./stdout)
- [Jaeger 导出器](./jaeger)
<!-- - [Zipkin 导出器](./zipkin) -->
- [OTLP-HTTP 导出器](./otlp-http)
<!-- - [OTLP-gRPC 导出器](./otlp-grpc) -->

## 使用方法

### 创建 Dubbo 实例

要使用链路追踪功能，你需要使用 `dubbo.NewInstance` 来创建一个 Dubbo 实例。

如果你直接创建 server 和 client，只能启用基本的 RPC 功能。

### 使用链路追踪功能

要配置链路追踪，你需要使用 `dubbo.WithTracing` 作为 `dubbo.NewInstance` 的选项。在 `dubbo.WithTracing` 中，你可以添加更多选项来配置链路追踪模块。

并且必须导入 `dubbo.apache.org/dubbo-go/v3/imports` 包。

```go
package main

import (
  "dubbo.apache.org/dubbo-go/v3"
  _ "dubbo.apache.org/dubbo-go/v3/imports"
  "dubbo.apache.org/dubbo-go/v3/otel/trace"
)

func main() {
    instance, err := dubbo.NewInstance(
        dubbo.WithTracing(
          // 在这里添加链路追踪选项
          trace.WithEnabled(), // 启用链路追踪功能
          trace.WithStdoutExporter(),
          trace.WithW3cPropagator(),
          trace.WithAlwaysMode(),
          trace.WithRatioMode(), // 使用比例模式
          trace.WithRatio(0.5), // 采样比例，仅在使用比例模式时生效
        ),
    )
}

```

如果你在 `dubbo.WithTracing` 中不添加任何选项，将使用默认的链路追踪配置。默认配置如下所示。

```yaml
# 默认链路追踪配置
enable: false
exporter: stdout
endpoint: ""
propagator: w3c
sample-mode: ratio
sample-ratio: 0.5
```

### 选项说明

- enable: 是否启用链路追踪
  - `trace.WithEnabled()` 表示启用链路追踪
- exporter: 链路追踪导出器后端，支持 stdout、jaeger、zipkin、otlp-http、otlp-grpc
  - `trace.WithStdoutExporter()`
  - `trace.WithJaegerExporter()`
  - `trace.WithZipkinExporter()`
  - `trace.WithOtlpHttpExporter()`
  - `trace.WithOtlpGrpcExporter()`
- endpoint: 导出器后端端点，例如，jaeger 导出器的端点是 `http://localhost:14268/api/traces`
  - `trace.WithEndpoint(string)`
- propagator: 上下文传播器类型，支持 w3c、b3，更多详情请参见[这里](https://opentelemetry.io/docs/concepts/context-propagation/)
  - `trace.WithW3cPropagator()` 
  - `trace.WithB3Propagator()` zipkin 导出器默认使用此选项
- sample-mode: 采样模式，支持 ratio（比例）、always（总是）、never（从不）
  - `trace.WithAlwaysMode()`
  - `trace.WithNeverMode()`
  - `trace.WithRatioMode()`
- sample-ratio: 采样比例，仅在采样模式为 ratio 时使用，范围在 0 到 1 之间
  - `trace.WithRatio(float64)`
