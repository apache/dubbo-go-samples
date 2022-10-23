# Dubbo-go Filter Samples

## Background

Dubbo-go designs and implements the filter mode, which helps the provider/consumer to process some actions before sending/processing requests. Dubbo-go has some built-in filter implementations, such as tps limiter, token, etc., and also supports user-defined implementations filter, please refer to the custom module.

## Instructions

### 1. Built-in Filter

#### 1.1 Background

Dubbo-go has built-in implementations of various filters, and the following types are connected by default during the running phase of the provider/consumer.
Provider：
- echo
- metrics
- token
- accesslog
- tps
- generic_service
- execute
- pshutdown

Consumer：
- cshutdown

For more built-in implementations of filter, please refer to the filter module of dubbo-go.

#### 1.2 Filter Confugration

Just configure the name of the filter under the corresponding provider or consumer, refer to Section 2.2.

### 2. Custom Filter

Take the consumer as an example, refer to the custom module for more specific code.

#### 2.1 Implement Custom Filter

Implement Filter interface.
```go
type Filter interface {
	Invoke(context.Context, protocol.Invoker, protocol.Invocation) protocol.Result
	OnResponse(context.Context, protocol.Result, protocol.Invoker, protocol.Invocation) protocol.Result
}
```
and inject it into the environment via `extension.SetFilter`：
```go
func init() {
	extension.SetFilter("myClientFilter", NewMyClientFilter)
}

func NewMyClientFilter() filter.Filter {
	return &MyClientFilter{}
}

type MyClientFilter struct {
}

func (f *MyClientFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	fmt.Println("MyClientFilter Invoke is called, method Name = ", invocation.MethodName())
	invocation.SetAttachment("request-key1", "request-value1")
	invocation.SetAttachment("request-key2", []string{"request-value2.1", "request-value2.2"})
	return invoker.Invoke(ctx, invocation)
}

func (f *MyClientFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, protocol protocol.Invocation) protocol.Result {
	fmt.Println("MyClientFilter OnResponse is called")
	fmt.Println("result attachment = ", result.Attachments())
	return result
}
```

#### 2.2 Filter Confugration

```yaml
# service config
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: 127.0.0.1:2181
  consumer:
    filter: myClientFilter
    check: true
    request_timeout: 3s
    connect_timeout: 3s
    references:
      GreeterClientImpl:
        protocol: tri
```

### 3. Running

See [HOWTO.md](../HOWTO_en.md) in the root directory to run this example.

Observe the output of the client:

```bash
MyClientFilter Invoke is called, method Name =  SayHello
MyClientFilter OnResponse is called
```