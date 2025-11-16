# 自定义 Filter 示例

[English](README.md) | [中文](README_zh.md)

## 背景

本示例演示了如何在 dubbo-go 中实现自定义过滤器。自定义过滤器允许您在客户端和服务端拦截和处理请求/响应。

## 实现方法

### 1. 实现自定义 Filter

实现 `Filter` 接口：

```go
type Filter interface {
	Invoke(context.Context, protocol.Invoker, protocol.Invocation) protocol.Result
	OnResponse(context.Context, protocol.Result, protocol.Invoker, protocol.Invocation) protocol.Result
}
```

### 2. 注册 Filter

在 `init()` 函数中使用 `extension.SetFilter()` 注册您的自定义过滤器：

**客户端 Filter** (`go-client/filter/myfilter.go`):
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

**服务端 Filter** (`go-server/filter/myfilter.go`):
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

### 3. 使用自定义 Filter

**客户端** (`go-client/cmd/main.go`):
```go
import (
	_ "github.com/apache/dubbo-go-samples/filter/custom/go-client/filter"
	"dubbo.apache.org/dubbo-go/v3/client"
	// ...
)

svc, err := greet.NewGreetService(cli, client.WithFilter("myClientFilter"))
```

**服务端** (`go-server/cmd/main.go`):
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

## 运行方法

### 前置条件

1. 启动 Zookeeper（默认：`127.0.0.1:2181`）

### 启动服务端

```shell
go run ./go-server/cmd/main.go
```

### 启动客户端

```shell
go run ./go-client/cmd/main.go
```

## 预期输出

**客户端输出:**
```bash
MyClientFilter Invoke is called, method Name =  Greet
MyClientFilter OnResponse is called
result attachment =  map[key1:value1 key2:[value1 value2]]
Greet response: hello world
```

**服务端输出:**
```bash
MyServerFilter Invoke is called, method Name =  Greet
request attachments = map[request-key1:request-value1 request-key2:[request-value2.1 request-value2.2]]
MyServerFilter OnResponse is called
```

## 注意事项

- 过滤器实现必须在单独的包中，以避免导入循环
- 使用 `_` 导入过滤器包以触发 `init()` 函数
- 过滤器可以在 `Invoke()` 中修改请求附件，在 `OnResponse()` 中修改响应附件

