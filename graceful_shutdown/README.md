# Graceful Shutdown Example

This sample is intended for manual integration testing of the graceful shutdown flow in `dubbo-go`.

It is useful for verifying these behaviors:
- active notice for long connections on Triple and gRPC
- passive closing behavior on the consumer side
- waiting for in-flight provider requests during shutdown
- the effect of `notify-timeout`, `step-timeout`, `consumer-update-wait`, and `offline-window`

This sample does **not** include a registry. That means you can test protocol-level active notice and request draining, but you cannot directly observe registry unregister propagation in this sample alone.

## Prerequisites

The sample already points to your local `E:/project/go_project/dubbo/dubbo-go` through `replace` in `go.mod`.

Run all commands from:

```bash
cd E:/project/go_project/dubbo/dubbo-go-samples/graceful_shutdown
```

## Quick Start

Start the server in one terminal:

```bash
go run ./go-server/cmd -protocols=tri -notify-timeout=5s -step-timeout=5s -delay=2s
```

Start the client in another terminal:

```bash
go run ./go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=3 -interval=300ms -request-timeout=6s
```

Then press `Ctrl+C` in the server terminal.

Expected result:
- in-flight requests may still complete
- new requests will begin to fail during shutdown
- server logs will show the graceful shutdown phases in order

## Important Address Format

For direct client calls, always include the protocol prefix in `-addr`.

Examples:
- Triple: `tri://127.0.0.1:20000`
- gRPC: `grpc://127.0.0.1:40000`
- Dubbo: `dubbo://127.0.0.1:30000`

If you omit the protocol prefix and only pass `127.0.0.1:20000`, the direct reference may be parsed incorrectly in some scenarios.

## Server Flags

`go-server/cmd/main.go` supports these test flags:
- `-protocols=tri|dubbo|grpc|tri,grpc`
- `-notify-timeout=5s`
- `-step-timeout=3s`
- `-consumer-update-wait=3s`
- `-offline-window=3s`
- `-delay=0s`

`-delay` adds artificial processing delay to every request so you can verify in-flight request draining.

## Client Flags

`go-client/cmd/main.go` supports these test flags:
- `-addr=tri://127.0.0.1:20000`
- `-interval=200ms`
- `-concurrency=1`
- `-request-timeout=5s`
- `-short=true|false`
- `-name-prefix=hello`

For long-connection testing, keep `-short=false`.

## Recommended Scenarios

### 1. Triple active notice with long connection

Terminal 1:
```bash
go run ./go-server/cmd -protocols=tri -notify-timeout=5s
```

Terminal 2:
```bash
go run ./go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=1 -interval=200ms
```

Then press `Ctrl+C` in the server terminal.

What to observe:
- server logs show the graceful shutdown phases
- client starts failing shortly after shutdown begins
- long connection is actively notified instead of only waiting for process exit

### 2. In-flight request draining

Terminal 1:
```bash
go run ./go-server/cmd -protocols=tri -delay=2s -step-timeout=5s
```

Terminal 2:
```bash
go run ./go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=3 -interval=300ms -request-timeout=6s
```

Then press `Ctrl+C` in the server terminal while requests are still running.

What to observe:
- already-running requests still have a chance to complete
- new requests begin to fail during shutdown
- the server exits after the in-flight wait budget is consumed or requests finish

### 3. Tight notify timeout

Terminal 1:
```bash
go run ./go-server/cmd -protocols=tri -notify-timeout=1s
```

Terminal 2:
```bash
go run ./go-client/cmd -addr=tri://127.0.0.1:20000
```

This is mainly for comparing server logs with a shorter active-notice budget.

### 4. Compare long and short connections

Long connection:
```bash
go run ./go-client/cmd -addr=tri://127.0.0.1:20000
```

Short connection:
```bash
go run ./go-client/cmd -addr=tri://127.0.0.1:20000 -short=true
```

Long connections are the more relevant path for active graceful notices.

### 5. Try gRPC manually

Terminal 1:
```bash
go run ./go-server/cmd -protocols=grpc -port-base=40000
```

Terminal 2:
```bash
go run ./go-client/cmd -addr=grpc://127.0.0.1:40000
```

Use this only if your local environment already has gRPC running stably in this sample.

## Practical Notes

- Triple is the best default protocol for manual verification.
- Dubbo can still help test request draining, but it does not cover the long-connection active notice path the same way.
- gRPC can also be tried if your environment is already stable for it.
- Because this sample has no registry, the "unregister from registry" phase is only part of the core implementation flow, not something you can fully observe here.