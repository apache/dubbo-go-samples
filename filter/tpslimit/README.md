# TPSLimit Filter Sample

### Background

Dubbo-go has a built in filter for limiting TPS purpose - "tpslimit". It can be enabled by configuring on the provider side, furthermore, user can customize the TPS limit strategy and the return value after the request is rejected.

### Example

##### 1. Code

A) Customize TPS limit strategy:

To customize TPS limit strategy, the interface "filter.TpsLimitStrategy" is needed to implement. In this example, the strategy is implemented as randomly rejecting the incoming request. Examples are linked as follows: https://github.com/apache/dubbo-go-samples/tree/master/filter/tpslimit/go-server/pkg/limit_strategy.go

B) Customize execution handler when the request is rejected.

Implement the interface "filter.RejectedExecutionHandler" to customize the return result to the client when the request is rejected. In this example, when the TPS limit criteria meets, the customized execution handler will return the error "The request is rejected and doesn't have any default value." back to the consumer. Examples are linked as follows: https://github.com/apache/dubbo-go-samples/tree/master/filter/tpslimit/go-server/pkg/reject_handler.go

##### 2. Configuration

Enable tpslimit filter in provider's configuration file like below:

```yaml
# service config
services:
  UserProvider:
    registry: demoZk
    protocol: dubbo
    interface: org.apache.dubbo.UserProvider
    tps.limiter: method-service
    tps.limit.strategy: RandomLimitStrategy
    tps.limit.rejected.handler: DefaultValueHandler
    tps.limit.interval: 5000
    tps.limit.rate: 300
```

##### 3. Run

Pls. refer to [HOWTO.md](../../HOWTO.md) under the root directory to run this sample.

The provider side will print out:

```bash
[2021-03-10/17:11:10 github.com/apache/dubbo-go-samples/filter/tpslimit/go-server/pkg.RandomTpsLimitStrategy.IsAllowable: limit_strategy.go: 56] %s
Random IsAllowable!
2021-03-10T17:11:10.748+0800 ERROR   filter_impl/tps_limit_filter.go:69      The invocation was rejected due to over the tps limitation, ...
```

The consumer side will print out:

```bash
error: The request is rejected and doesn't have any default value. 
```