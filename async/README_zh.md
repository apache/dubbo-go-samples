# Async 示例

### 背景

Dubbo-go 提供同步调用的同时，还提供了客户端异步调用的能力。客户端可以通过实现以下的接口来异步获得服务端返回的响应：


```golang
type AsyncCallbackService interface {
	CallBack(response CallbackResponse)
}
```

### 示例

**代码**

```golang
type UserProvider struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *User) error
}

func (u *UserProvider) CallBack(res common.CallbackResponse) {
	fmt.Println("CallBack res:",res)
}
```

**配置**

提供回调方法的同时，还需要在客户端的配置中增加 **"async:true"** 的配置，如下所示：

```yaml
# reference config
references:
  "UserProvider":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.UserProvider"
    # this is necessary to enable async call
    async: true
```

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。


