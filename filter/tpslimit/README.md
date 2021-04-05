# TPSLimit Filter Sample

### Background

Dubbo-go has a builtin filter for limiting TPS purpose - "tpslimit". It can be enabled by configuring on the provider side, furthermore, user can customize the TPS limit strategy and the return value after the request is rejected.

### Example

##### 1. Code

A) Customize TPS limit strategy:

To customize TPS limit strategy, the interface "filter.TpsLimitStrategy" is needed to implement. In this example, the strategy is implemented as randomly rejecting the incoming request, see below: 

```go
func init() {
	/*
	 * register your implementation and them using it like:
	 *
	 * "UserProvider":
	 *   registry: "hangzhouzk"
	 *   protocol : "dubbo"
	 *   interface : "com.ikurento.user.UserProvider"
	 *   ... # other configuration
	 *   tps.limiter: "method-service" # the name of limiter
	 *   tps.limit.strategy: "RandomLimitStrategy"
	 */
	extension.SetTpsLimitStrategy("RandomLimitStrategy", &RandomTpsLimitStrategyCreator{})
}

/**
 * The RandomTpsLimitStrategy should not be singleton because different TpsLimiter will create many instances.
 * we won't want them affect each other.
 */
type RandomTpsLimitStrategy struct {
	rate     int
	interval int
}

func (r RandomTpsLimitStrategy) IsAllowable() bool {
	// this is a simple demo.
	gxlog.CInfo("Random IsAllowable!")
	randNum := rand.Int63n(2)
	return randNum == 0
}

type RandomTpsLimitStrategyCreator struct{}

func (creator *RandomTpsLimitStrategyCreator) Create(rate int, interval int) filter.TpsLimitStrategy {
	return &RandomTpsLimitStrategy{
		rate:     rate,
		interval: interval,
	}
}
```

B) Customize execution handler when the request is rejected.

Implement the interface "filter.RejectedExecutionHandler" to customize the return result to the client when the request is rejected. In this example, when the TPS limit criteria meets, the customized execution handler will return the error "The request is rejected and doesn't have any default value." back to the consumer.

```go
func init() {
	/*
	 * register your custom implementation into filter.
	 * "DefaultValueHandler" is the name used in configure file, like server.yml:
	 * "UserProvider":
	 *   registry: "hangzhouzk"
	 *   protocol : "dubbo"
	 *   interface : "com.ikurento.user.UserProvider"
	 *   ... # other configuration
	 *   tps.limiter: "method-service",
	 *
	 *   tps.limit.rejected.handler: "DefaultValueHandler",
	 * So when the invocation is over the tps limitation, it will return the default value.
	 * This is a common use case.
	 */
	extension.SetRejectedExecutionHandler("DefaultValueHandler", GetDefaultValueRejectedExecutionHandlerSingleton)

}

/**
 * The RejectedExecutionHandler is used by some components,
 * e.g, ExecuteLimitFilter, GracefulShutdownFilter, TpsLimitFilter.
 * When the requests are rejected, the RejectedExecutionHandler allows you to do something.
 * You can alert the developer, or redirect those requests to another providers. It depends on what you need.
 *
 * Let's assume that you need a RejectedExecutionHandler which will return some default result if the request was rejected.
 */
type DefaultValueRejectedExecutionHandler struct {
	defaultResult sync.Map
}

func (mh *DefaultValueRejectedExecutionHandler) RejectedExecution(url *common.URL, invocation protocol.Invocation) protocol.Result {
	// put your custom business here.
	logger.Error("Here is my custom rejected handler. I want to do something if the requests are rejected. ")
	// in most cases, if the request was rejected, you won't want to invoke the origin provider.
	// But if you really want to do that, you can do it like this:
	// invocation.Invoker().Invoke(invocation)

	// the ServiceKey + methodName is the key
	key := url.ServiceKey() + "#" + invocation.MethodName()
	result, loaded := mh.defaultResult.Load(key)
	if !loaded {
		// we didn't configure any default value for this invocation
		return &protocol.RPCResult{
			Err: errors.New("The request is rejected and doesn't have any default value. "),
		}
	}
	return result.(*protocol.RPCResult)
}

func GetCustomRejectedExecutionHandler() filter.RejectedExecutionHandler {
	return &DefaultValueRejectedExecutionHandler{}
}

var (
	customHandlerOnce     sync.Once
	customHandlerInstance *DefaultValueRejectedExecutionHandler
)

/**
 * the better way is designing the RejectedExecutionHandler as singleton.
 */
func GetDefaultValueRejectedExecutionHandlerSingleton() filter.RejectedExecutionHandler {
	customHandlerOnce.Do(func() {
		customHandlerInstance = &DefaultValueRejectedExecutionHandler{}
	})

	initDefaultValue()

	return customHandlerInstance
}

func initDefaultValue() {
	// setting your default value
}
```

##### 2. Configuration

Enable tpslimit filter in provider's configuration file like below:

```yaml
# service config
services:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    tps.limiter: "method-service"
    tps.limit.strategy: "RandomLimitStrategy"
    tps.limit.rejected.handler: "DefaultValueHandler"
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