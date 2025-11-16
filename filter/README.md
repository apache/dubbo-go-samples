# Dubbo-go Filter Samples

[English](README.md) | [中文](README_zh.md)

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

#### 1.2 Filter Configuration

Use `server.WithFilter()` for provider side or `client.WithFilter()` for consumer side. See examples in the token, sentinel, and custom modules.

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

#### 2.2 Filter Configuration

For provider side:
```go
if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{},
	server.WithFilter("myServerFilter"),
); err != nil {
	panic(err)
}
```

For consumer side:
```go
svc, err := greet.NewGreetService(cli, client.WithFilter("myClientFilter"))
if err != nil {
	panic(err)
}
```

## Examples

This directory contains the following filter examples:

- **custom**: Custom filter implementation example (client and server)
- **token**: Token filter example
- **sentinel**: Sentinel filter example for flow control and circuit breaking
- **polaris/limit**: Polaris TPS limiter example
- **tpslimit**: Custom TPS limit strategy and rejected execution handler example

