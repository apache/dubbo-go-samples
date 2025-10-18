# otlp http tracing exporter

English | [中文](README_CN.md)

The integration checks dubbo/triple/jsonrpc tracing feature, and uses OpenTelemetry's otlp http exporter to export the tracing data.

The test exports the tracing data to a local mock http server.

## How to run

### Run server

```shell
go run ./go-server/cmd/main.go
```

### Run test (client)

```shell
go test -tags integration -v ./tests/integration/... 
```

If test success, you will see the log like below in the server terminal:

```shell
2025-09-18 16:22:29	INFO	cmd/main.go:127	server count: 3, client count: 3
```

If the test fails, you will see the panic in the server terminal.

