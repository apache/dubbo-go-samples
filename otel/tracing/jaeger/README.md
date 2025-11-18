# Jaeger tracing exporter

[English](README.md) | [中文](README_zh.md)

This example shows dubbo-go's tracing feature with Jaeger exporter.

## Prerequisites

Before running this example, you need to start a Jaeger instance. You can use Docker to run Jaeger:

```shell
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 4317:4317 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

Or use docker-compose:

```yaml
version: '3'
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # Jaeger UI
      - "14268:14268"  # Jaeger collector HTTP endpoint
    environment:
      - COLLECTOR_OTLP_ENABLED=true
```

## How to run

### Run server

```shell
$ go run ./go-server/cmd/main.go
```

### Run client

```shell
$ go run ./go-client/cmd/main.go
```

## View traces

After running the client and server, you can view the traces in the Jaeger UI:

1. Open your browser and navigate to `http://localhost:16686`
2. Select the service name (e.g., `dubbo_otel_jaeger_server` or `dubbo_otel_jaeger_client`)
3. Click "Find Traces" to see the tracing data

## Configuration

The Jaeger exporter is configured with the following options:

- `trace.WithJaegerExporter()`: Use Jaeger as the tracing exporter
- `trace.WithEndpoint("http://localhost:14268/api/traces")`: Set the Jaeger collector endpoint
- `trace.WithW3cPropagator()`: Use W3C trace context propagation
- `trace.WithAlwaysMode()`: Always sample traces

You can modify the endpoint to point to your Jaeger instance if it's running on a different host or port.

