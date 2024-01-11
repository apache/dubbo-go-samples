# Filter 示例

## 背景

Dubbo-go 设计实现了过滤器模式，帮助客户端/服务端在发送/处理请求前置处理一些动作，dubbo-go 内置了一些过滤器实现，如 tps limiter、token 等，也支持用户自定义实现 filter，可参考 custom 模块。

## 使用方法

### 1. 内置 Filter

#### 1.1 概览
Dubbo-go 内置实现了多种 filter，并在框架运行阶段默认接入了以下几种：
Provider 端：
- echo
- metrics
- token
- accesslog
- tps
- generic_service
- execute
- pshutdown  

Consumer 端：
- cshutdown  

更多 filter 内置实现可参考 dubbo-go 的 filter 模块。

#### 1.2 Filter 配置
只需在对应 provider 或 consumer 下配置该 filter 的名字即可，参考 2.2 节。

### 2. 自定义 Filter
以 Client 端为例，更具体的代码参考 custom 模块。

#### 2.1 实现自定义的 Filter

实现 Filter 接口
```go
type Filter interface {
	Invoke(context.Context, protocol.Invoker, protocol.Invocation) protocol.Result
	OnResponse(context.Context, protocol.Result, protocol.Invoker, protocol.Invocation) protocol.Result
}
```
并通过 `extension.SetFilter` 将其注入到环境中，如 custom 中：
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

#### 2.2 Filter 配置

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

### 3. 运行

请参阅根目录中的 [HOWTO.md](../../HOWTO_zh.md) 来运行本例。

观察客户端的输出：

```bash
MyClientFilter Invoke is called, method Name =  SayHello
MyClientFilter OnResponse is called
```