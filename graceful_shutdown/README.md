# Graceful Shutdown Example

English | [中文](README_CN.md)

This sample is intended for manual verification of the Triple graceful shutdown flow in `dubbo-go`.

It is useful for verifying these behaviors:
- active notice for long connections on Triple
- passive closing behavior on the consumer side
- waiting for in-flight provider requests during shutdown
- the effect of `timeout`, `step-timeout`, `consumer-update-wait`, and `offline-window`

This sample does **not** include a registry. That means you can test protocol-level active notice and request draining, but you cannot directly observe registry unregister propagation in this sample alone.

## Prerequisites

This sample uses the repository root `go.mod`.

Run all commands from your local checkout of `dubbo-go-samples`:

```bash
cd /path/to/dubbo-go-samples
```

## Quick Start

Start the server in one terminal:

```bash
go run ./graceful_shutdown/go-server/cmd -timeout=60s -step-timeout=5s -delay=2s
```

Start the client in another terminal:

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=3 -interval=300ms -request-timeout=6s
```

Then press `Ctrl+C` in the server terminal.

Expected result:
- in-flight requests may still complete
- new requests will begin to fail during shutdown
- server logs will show the graceful shutdown phases in order

## Important Address Format

For direct client calls, always include the protocol prefix in `-addr`.

Example:
- Triple: `tri://127.0.0.1:20000`

If you omit the protocol prefix and only pass `127.0.0.1:20000`, the direct reference may be parsed incorrectly in some scenarios.

## Server Flags

`graceful_shutdown/go-server/cmd/main.go` supports these test flags:
- `-port=20000`
- `-timeout=60s`
- `-step-timeout=3s`
- `-consumer-update-wait=3s`
- `-offline-window=3s`
- `-delay=0s`

`-delay` adds artificial processing delay to every request so you can verify in-flight request draining.

## Client Flags

`graceful_shutdown/go-client/cmd/main.go` supports these test flags:
- `-addr=tri://127.0.0.1:20000`
- `-interval=200ms`
- `-concurrency=1`
- `-request-timeout=5s`
- `-short=true|false`
- `-name-prefix=hello`
- `-max-requests=0`
- `-min-successes=0`
- `-min-failures=0`

For long-connection testing, keep `-short=false`.

`-max-requests`, `-min-successes`, and `-min-failures` are mainly for automated verification. The client panics if the configured minimum counts are not reached before exit.

## Recommended Scenarios

### 1. Triple active notice with long connection

Terminal 1:
```bash
go run ./graceful_shutdown/go-server/cmd -timeout=60s
```

Terminal 2:
```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=1 -interval=200ms
```

Then press `Ctrl+C` in the server terminal.

What to observe:
- server logs show the graceful shutdown phases
- client starts failing shortly after shutdown begins
- long connection is actively notified instead of only waiting for process exit

### 2. In-flight request draining

Terminal 1:
```bash
go run ./graceful_shutdown/go-server/cmd -delay=2s -step-timeout=5s
```

Terminal 2:
```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=3 -interval=300ms -request-timeout=6s
```

Then press `Ctrl+C` in the server terminal while requests are still running.

What to observe:
- already-running requests still have a chance to complete
- new requests begin to fail during shutdown
- the server exits after the in-flight wait budget is consumed or requests finish

### 3. Observe active notice and request draining together

Use a short consumer update wait so shutdown starts rejecting new work earlier while existing requests are still draining.

Terminal 1:
```bash
go run ./graceful_shutdown/go-server/cmd -delay=2s -timeout=15s -step-timeout=2s -consumer-update-wait=0s
```

Terminal 2:
```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=2 -interval=200ms -request-timeout=4s
```

Then press `Ctrl+C` in the server terminal.

What to observe:
- server logs print the full graceful shutdown sequence
- some in-flight requests still complete after shutdown starts
- newer requests begin to fail earlier than in the default configuration
- client logs include the Triple active-notice path from `triple-health-watch`

### 4. Tight overall timeout

Terminal 1:
```bash
go run ./graceful_shutdown/go-server/cmd -timeout=10s -step-timeout=1s
```

Terminal 2:
```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000
```

This is mainly for comparing server logs with a tighter overall graceful shutdown budget.

### 5. Compare long and short connections

Long connection:
```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000
```

Short connection:
```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -short=true
```

Long connections are the more relevant path for active graceful notices.

## Integration Test

This sample is wired into the root integration test flow:

```bash
./integrate_test.sh graceful_shutdown
```

The script starts the Triple server, runs the client in the background, waits until at least one request succeeds, and then sends an interrupt signal to trigger graceful shutdown.

Before the client exits, it must observe:

- at least one successful request
- at least one failed request during shutdown

If those expectations are not met, the client panics so CI fails immediately.

## Practical Notes

- Triple is the intended protocol for manual verification in this sample.
- This sample is intentionally Triple-only so it focuses on the active notice path implemented in the current graceful shutdown flow.
- Because this sample has no registry, the "unregister from registry" phase is only part of the core implementation flow, not something you can fully observe here.
