# Custom Filter Sample

[English](README.md) | [中文](README_zh.md)

## Background

This example demonstrates how to implement custom filters in dubbo-go. Custom filters allow you to intercept and process requests/responses on both the client and server sides.

## Implementation

### 1. Implement Custom Filter

Implement the `Filter` interface:

```go
type Filter interface {
	Invoke(context.Context, protocol.Invoker, protocol.Invocation) protocol.Result
	OnResponse(context.Context, protocol.Result, protocol.Invoker, protocol.Invocation) protocol.Result
}
```

### 2. Register Filter

Register your custom filter using `extension.SetFilter()` in the `init()` function:

**Client Filter** (`go-client/filter/myfilter.go`):
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

**Server Filter** (`go-server/filter/myfilter.go`):
```go
func init() {
	extension.SetFilter("myServerFilter", NewMyServerFilter)
}

func NewMyServerFilter() filter.Filter {
	return &MyServerFilter{}
}

type MyServerFilter struct {
}

func (f *MyServerFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
	fmt.Println("MyServerFilter Invoke is called, method Name = ", invocation.MethodName())
	fmt.Printf("request attachments = %s\n", invocation.Attachments())
	return invoker.Invoke(ctx, invocation)
}

func (f *MyServerFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, protocol protocol.Invocation) protocol.Result {
	fmt.Println("MyServerFilter OnResponse is called")
	myAttachmentMap := make(map[string]interface{})
	myAttachmentMap["key1"] = "value1"
	myAttachmentMap["key2"] = []string{"value1", "value2"}
	result.SetAttachments(myAttachmentMap)
	return result
}
```

### 3. Use Custom Filter

**Client Side** (`go-client/cmd/main.go`):
```go
import (
	_ "github.com/apache/dubbo-go-samples/filter/custom/go-client/filter"
	"dubbo.apache.org/dubbo-go/v3/client"
	// ...
)

svc, err := greet.NewGreetService(cli, client.WithFilter("myClientFilter"))
```

**Server Side** (`go-server/cmd/main.go`):
```go
import (
	_ "github.com/apache/dubbo-go-samples/filter/custom/go-server/filter"
	"dubbo.apache.org/dubbo-go/v3/server"
	// ...
)

if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{},
	server.WithFilter("myServerFilter"),
); err != nil {
	panic(err)
}
```

## Running

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

**Client Output:**
```bash
MyClientFilter Invoke is called, method Name =  Greet
MyClientFilter OnResponse is called
result attachment =  map[key1:value1 key2:[value1 value2]]
Greet response: hello world
```

**Server Output:**
```bash
MyServerFilter Invoke is called, method Name =  Greet
request attachments = map[request-key1:request-value1 request-key2:[request-value2.1 request-value2.2]]
MyServerFilter OnResponse is called
```

## Notes

- The filter implementation must be in a separate package to avoid import cycles
- Import the filter package with `_` to trigger the `init()` function
- Filters can modify request attachments in `Invoke()` and response attachments in `OnResponse()`

