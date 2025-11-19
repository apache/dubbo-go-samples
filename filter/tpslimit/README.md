# TPSLimit Filter Sample

[English](README.md) | [中文](README_zh.md)

### Background

Dubbo-go has a built in filter for limiting TPS purpose - "tpslimit". It can be enabled by configuring on the provider side, furthermore, user can customize the TPS limit strategy and the return value after the request is rejected.

### Example

##### 1. Code

A) Customize TPS limit strategy:

To customize TPS limit strategy, the interface "filter.TpsLimitStrategy" is needed to implement. In this example, the strategy is implemented as randomly rejecting the incoming request. Examples are linked as follows: [limit_strategy.go](go-server/pkg/limit_strategy.go)

B) Customize execution handler when the request is rejected.

Implement the interface "filter.RejectedExecutionHandler" to customize the return result to the client when the request is rejected. In this example, when the TPS limit criteria meets, the customized execution handler will return the error "The request is rejected and doesn't have any default value." back to the consumer. Examples are linked as follows: [reject_handler.go](go-server/pkg/reject_handler.go)

##### 2. Configuration

Enable tpslimit filter in provider's code using v3 API:

```go
import "dubbo.apache.org/dubbo-go/v3/config"

if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{},
	server.WithTpsLimiter("method-service"),
	server.WithMethod(
		config.WithName("Greet"),
		config.WithTpsLimitRate(5),        // must be >0, otherwise the limiter falls back to the default -1 and gets ignored
		config.WithTpsLimitInterval(1000), // ms
		config.WithTpsLimitStrategy("RandomLimitStrategy"),
	),
	server.WithTpsLimitRejectedHandler("DefaultValueHandler"),
); err != nil {
	panic(err)
}
```

> **Note**: Only configuring `tps.limit.rate` at the service level may be overridden by the Provider defaults (which set `-1`). Declaring the TPS options via `server.WithMethod` guarantees `greet.GreetService#Greet` receives a positive rate and interval so that the limiter actually works.

The custom TPS limit strategy and rejected execution handler are registered via `extension.SetTpsLimitStrategy()` and `extension.SetRejectedExecutionHandler()` in the `init()` function of the respective packages.

##### 3. Run

### Prerequisites

1. Start Zookeeper (default: `127.0.0.1:2181`)

### Run Server

```shell
go run ./go-server/cmd/main.go
```

### Run Client

```shell
go run ./go-client/cmd/main.go
```

## Expected Output

**Server Output:**

```bash
Random IsAllowable!
ERROR   The invocation was rejected due to over the tps limitation, ...
```

**Client Output:**

```bash
start to test tpslimit
error: The request is rejected and doesn't have any default value.
response: hello world
...
successCount=<number>, failCount=<number>
```

The client will send 60 requests with 200ms intervals. Some requests will be rejected by the TPS limiter, and you'll see both success and failure counts in the final output.

