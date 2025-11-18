# Tracing

English | [中文](./README_zh.md)

## Background

Dubbo-go supports [OpenTelemetry](https://opentelemetry.io/) tracing.

- [Stdout exporter](./stdout)
- [Jaeger exporter](./jaeger)
<!-- - [Zipkin exporter](./zipkin) -->
- [OTLP-HTTP exporter](./otlp-http)
<!-- - [OTLP-gRPC exporter](./otlp-grpc) -->

## Usage

### Create Dubbo instance

To use the tracing feature, you need to use `dubbo.NewInstance` to create a Dubbo instance. 

If you create server and client directly, only basic RPC feature can be enabled.

### Use tracing feature

To config tracing, you need to use `dubbo.WithTracing` as an option of `dubbo.NewInstance`. In `dubbo.WithTracing`, you can add more options to config tracing module.

And the `dubbo.apache.org/dubbo-go/v3/imports` must be imported.

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
          // add tracing options here
          trace.WithEnabled(), // enable tracing feature
          trace.WithStdoutExporter(),
          trace.WithW3cPropagator(),
          trace.WithAlwaysMode(),
          trace.WithRatioMode(), // use ratio mode
          trace.WithRatio(0.5), // sample ratio, only active when use ratio mode
        ),
    )
}

```

If you don't add any options in `dubbo.WithTracing`, default tracing config will be used. Default tracing config is shown below.

```yaml
# default tracing config
enable: false
exporter: stdout
endpoint: ""
propagator: w3c
sample-mode: ratio
sample-ratio: 0.5
```

### Options description

- enable: enable tracing or not
  - `trace.WithEnabled()` means enable tracing
- exporter: tracing exporter backends, support stdout, jaeger, zipkin, otlp-http, otlp-grpc
  - `trace.WithStdoutExporter()`
  - `trace.WithJaegerExporter()`
  - `trace.WithZipkinExporter()`
  - `trace.WithOtlpHttpExporter()`
  - `trace.WithOtlpGrpcExporter()`
- endpoint: exporter backend endpoint, for example, jaeger exporter's endpoint is `http://localhost:14268/api/traces`
  - `trace.WithEndpoint(string)`
- propagator: context propagator type, support w3c, b3, more details you can see [here](https://opentelemetry.io/docs/concepts/context-propagation/)
  - `trace.WithW3cPropagator()` 
  - `trace.WithB3Propagator()` zipkin exporter default use this
- sample-mode: sample mode, support ratio, always, never
  - `trace.WithAlwaysMode()`
  - `trace.WithNeverMode()`
  - `trace.WithRatioMode()`
- sample-ratio: sample ratio, only used when sample-mode is ratio, range between 0 and 1
  - `trace.WithRatio(float64)`


