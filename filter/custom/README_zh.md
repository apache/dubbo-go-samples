# 自定义 Filter 示例

### 背景

Dubbo-go 支持在服务端和客户端扩展自己的 filter，通过这样做，在请求链路上做一些自定义的处理，比如，监控或者审计等。

### 示例

在本例中，将以服务端为例，展示如何增加一个自定义的 filter。

##### 1. 代码

A) 首先在服务端实现一个自定义的 filter，如下所示：

```go
type myCustomFilter struct{}

func (mf myCustomFilter) Invoke(ctx context.Context, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
    // the logic put here...
    // you can get many params in url. And the invocation provides more information about
    url := invoker.GetURL()
    serviceKey := url.ServiceKey()
    gxlog.CInfo("Here is the my custom filter. The service is invoked: %s", serviceKey)
    return invoker.Invoke(ctx, invocation)
}

func (mf myCustomFilter) OnResponse(ctx context.Context, result protocol.Result, invoker protocol.Invoker, invocation protocol.Invocation) protocol.Result {
    // you can do something here with result
    gxlog.CInfo("Got result!")
    if user, ok := result.Result().(*User); ok {
        user.Name = strings.ToUpper(user.Name)
        user.Age = user.Age + 10
    }
    return result
}

func GetMyCustomFilter() filter.Filter {
    return &myCustomFilter{}
}
```

该 filter 会修改返回给客户端的结果，将用户名改为大写，并将用户的年龄加十。同时，还将在服务端打印日志。

B) 把该 filter 注册成名为 "MyCustomFilter" 的 filter：

```go
func init() {
	/**
	 * MyCustomFilter would be the name that used in your configuration file.
	 * it can be used as reference filter and provider filter.
	 * For example, using this filter in server, and the configure file looks like:
	 *
	 * filter: "MyCustomFilter",
	 * registries:
	 *  "demoZk":
	 *    protocol: "zookeeper"
	 *    timeout	: "3s"
	 *    address: "127.0.0.1:2181"
	 * Another important things is that you should make sure this statement executed. It usually means that
	 * this file should be imported.
	 */
	extension.SetFilter("MyCustomFilter", GetMyCustomFilter)

	// or using the singleton
	// filter.SetFilter("MyCustomFilter", GetMyCustomFilterSingleton)
}
```

##### 2. 配置

在服务端的配置文件中，按如下所示配置该 filter：

```yaml
# filter config
filter: "MyCustomFilter"
```

##### 3. 运行

请参阅根目录中的 [HOWTO.md](../../HOWTO_zh.md) 来运行本例。

观察服务端的输出：

```bash
[2021-03-10/16:30:52 github.com/apache/dubbo-go-samples/filter/custom/go-server/pkg.myCustomFilter.Invoke: custom_filter.go: 61] %s
Here is the my custom filter. The service is invoked: org.apache.dubbo.UserProvider
```

观察客户端的输出：

```bash
[2021-03-10/16:32:06 main.main: client.go: 64] %s response result: &{A001 ALEX STOCKS 28 2021-03-10 16:32:06.643 +0800 CST}
```