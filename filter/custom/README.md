# Custom Filter Sample

### Background

Dubbo-go supports custom filter on both provider side and consumer side. By doing this, user has the opportunity to do some customized operations during the request, for example: logging or auditing, etc.

### Example

This example shows how to add a custom filter on the provider side.

##### 1. Code

A) First implement a custom filter on the provider side as shown below:

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

This filter implementation will modify the response result before send it back to the consumer, for example: capitalize user's name, and add user's age by 10. At the same time, this filter also prints out messages for logging purpose.

B) Register this custom filter with the name "MyCustomFilter":

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

##### 2. Configuration

Configure this filter in the provider's configuration file like this:

```yaml
# filter config
filter: "MyCustomFilter"
```

##### 3. Run

Pls. refer to [HOWTO.md](../../HOWTO.md) under the root directory to run this sample.

The provider side will print out:

```bash
[2021-03-10/16:30:52 github.com/apache/dubbo-go-samples/filter/custom/go-server/pkg.myCustomFilter.Invoke: custom_filter.go: 61] %s
Here is the my custom filter. The service is invoked: org.apache.dubbo.UserProvider
```

And the consumer side will print out:

```bash
[2021-03-10/16:32:06 main.main: client.go: 64] %s response result: &{A001 ALEX STOCKS 28 2021-03-10 16:32:06.643 +0800 CST}
``
