# Async Sample

### Background

Dubbo-go provides not only synchronous invocation, but also asynchronous invocation on the consumer side. In order to
use it, the client needs to implement the following interface to asynchronously receive the response from the service
provider:

```golang
type AsyncCallbackService interface {
CallBack(response CallbackResponse)
}
```

### Example

**Code**

```golang
type UserProvider struct {
    GetUser func (ctx context.Context, req []interface{}, rsp *User) error
}

func (u *UserProvider) CallBack(res common.CallbackResponse) {
    fmt.Println("CallBack res:", res)
}
```

**Configuration**

Besides, client also needs to config **"async:true"** in consumer's yaml config file as following:

```yaml
# reference config
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    # this is necessary to enable async call
    async: true
```

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.


