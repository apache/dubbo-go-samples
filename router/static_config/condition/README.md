# Static condition router

This example shows how to configure Dubbo-Go's condition router statically in code with direct provider URLs.

English | [中文](README_CN.md)

## Prerequisites

- Go 1.25+.

## What this sample demonstrates

- Two providers listening directly on ports `20000` and `20001`
- A service-scope static condition router configured by `dubbo.WithRouter(...)`
- A `force=true` rule that routes `Greet` requests to port `20000`

## How to run

### Run providers

Start two providers in separate terminals:

```shell
$ go run ./go-server/cmd/server.go
$ go run ./go-node2-server/cmd/server_node2.go
```

- `go-server` listens on `:20000`
- `go-node2-server` listens on `:20001`

### Run consumer

If you want to change the provider addresses, update `directURL` in `go-client/cmd/client.go`.

```shell
$ go run ./go-client/cmd/client.go
```

The client connects to both providers by direct URLs.
No registry or config center is required.

## Expected result

- The client logs `invoke successfully: receive: static condition router, response from: server-node-20000`

## Key router config

The consumer injects this service-scope static condition router:

```go
dubbo.WithRouter(
    router.WithScope("service"),
    router.WithKey(greet.GreetServiceName),
    router.WithForce(true),
    router.WithConditions([]string{
        "method = Greet => port = 20000",
    }),
)
```
