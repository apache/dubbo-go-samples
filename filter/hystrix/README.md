# Hystrix Filter Example

[English](README.md) | [中文](README_zh.md)

## Background

This example demonstrates how to use the Hystrix filter in dubbo-go to implement circuit breaker functionality. Hystrix is a latency and fault tolerance library designed to isolate points of access to remote systems, services, and 3rd party libraries, stop cascading failures, and enable resilience in complex distributed systems.

## Implementation

### 1. Configure Hystrix Commands

Use the `hystrix-go` API to configure circuit breaker commands. The resource name format is:
```
dubbo:consumer:InterfaceName:group:version:Method
```

**Client Configuration** (`go-client/cmd/main.go`):
```go
import (
    "github.com/afex/hystrix-go/hystrix"
    _ "github.com/apache/dubbo-go-extensions/filter/hystrix"
)

func init() {
    // Resource name format: dubbo:consumer:InterfaceName:group:version:Method
    cmdName := "dubbo:consumer:greet.GreetService:::Greet"

    hystrix.ConfigureCommand(cmdName, hystrix.CommandConfig{
        Timeout:                1000, // timeout in milliseconds
        MaxConcurrentRequests:  10,   // max concurrent requests
        RequestVolumeThreshold: 5,    // minimum requests to trip the circuit
        SleepWindow:            5000,  // time to wait before attempting recovery (ms)
        ErrorPercentThreshold:  50,   // error rate threshold (percentage)
    })
}
```

### 2. Use Hystrix Filter

**Client** (`go-client/cmd/main.go`):
```go
import (
    "dubbo.apache.org/dubbo-go/v3/client"
)

svc, err := greet.NewGreetService(cli, client.WithFilter("hystrix_consumer"))
```

## Configuration Parameters

| Parameter | Description |
|-----------|-------------|
| `Timeout` | Command execution timeout in milliseconds |
| `MaxConcurrentRequests` | Maximum number of concurrent requests allowed |
| `RequestVolumeThreshold` | Minimum number of requests required to trip the circuit (within sliding window) |
| `SleepWindow` | Time to wait after circuit opens before attempting recovery (milliseconds) |
| `ErrorPercentThreshold` | Error rate threshold that triggers circuit opening (percentage) |

## How to Run

### Prerequisites

1. Start the Go server and make sure `127.0.0.1:20000` is ready to accept requests.

### Start Server

```shell
cd filter/hystrix/go-server
go run ./cmd/main.go
```

### Start Client

```shell
cd filter/hystrix/go-client
go run ./cmd/main.go
```

## Expected Output

**Client Output:**
```bash
=== Test 1: Sending normal requests ===
Request 1 success: Hello, request-1! (request #1)
Request 2 success: Hello, request-2! (request #2)
Request 3 success: Hello, request-3! (request #3)

=== Test 2: Sending concurrent requests ===
Concurrent request 1 success: Hello, concurrent-1! (request #4)
Concurrent request 2 success: Hello, concurrent-2! (request #5)
...

=== Test 3: Sending requests after concurrent test ===
After-test request 1 failed (circuit might be open): ...
```

When the circuit breaker is open, you will see errors like:
```bash
After-test request 1 failed: hystrix: circuit open
```

## Testing Circuit Breaker

The example program includes three test phases:

1. **Normal Requests**: Sends 3 normal requests to verify basic functionality
2. **Concurrent Requests**: Sends 15 concurrent requests, potentially triggering the circuit breaker
3. **Recovery Test**: Continues sending requests after concurrent test to observe circuit breaker state

If the circuit breaker is triggered, wait about 5 seconds (SleepWindow configuration) before running the client again to see the circuit recover.

## Notes

- Hystrix filter is primarily used on the **client side** to protect callers from downstream service failures
- The resource name must match the actual interface, group, version, and method name
- Circuit breaker states: Closed → Open → Half-Open → Closed
- Configure timeout and concurrency limits appropriately to avoid resource exhaustion
