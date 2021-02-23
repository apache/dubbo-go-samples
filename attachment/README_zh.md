# Attachment 示例

### 背景

可以通过 attachment 把用户的数据从 Dubbo 的客户端传递给服务端。在 Dubbo-go 中，attachment 在二者之间的传递是通过 `context.Context` 的机制来完成的。如果要使用 attachment，那么应当把 `context.Context` 作为要调用的服务方法的第一个参数，例如:

```go
GetUser func(ctx context.Context, req []interface{}, rsp *User) error
```

为了在客户端传递用户数据到服务端，首先需要在方法的第一个参数 `context.Context` 中放入一个 `map[string]interface{}` 的数据类型，该数据的 Key 值约定为 "attachment"。以下的代码片段展示了如何把用户自定义的的一个时间戳通过该 map 中的 "timestamp" 放入 attachment 中:

```go
ctx := context.WithValue(context.Background(), constant.AttachmentKey, 
	map[string]interface{}{"timestamp": time.Now()})
err := userProvider.GetUser(ctx, []interface{}{"A001"}, user)
```

在服务提供方，方法的第一个参数 `context.Context` 传入时携带了用户在客户端放入的自定义数据。下面的代码就是本例子中的服务端实现，主要展示了如何从 attachment 中提取用户自定义的数据:


```go
func (u *UserProvider) GetUser(ctx context.Context, req []interface{}) (*User, error) {
	t := time.Now()
	attachment := ctx.Value(constant.AttachmentKey).(map[string]interface{})
	if v, ok := attachment["timestamp"]; ok {
		t = v.(time.Time).Add(-1 * 365 * 24 * time.Hour)
	}

	rsp := User{"A001", "Alex Stocks", 18, t}
	return &rsp, nil
}
```

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。


