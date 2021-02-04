# Attachment Example

### Background

A Dubbo client can pass user data to the remote Dubbo server via attachment. Dubbo-go leverages `context.Context` as the attachment between consumer and provider. In order to use attachment, a `context.Context` is required to introduce as the first parameter of the service method, for example:

```go
GetUser func(ctx context.Context, req []interface{}, rsp *User) error
```

To pass the user data from the client side, a `map[string]interface{}` should be put into the `context.Context` with the key "attachment". The following code snippet shows how a user data "timestamp" is put into the attachment:

```go
ctx := context.WithValue(context.Background(), constant.AttachmentKey, 
	map[string]interface{}{"timestamp": time.Now()})
err := userProvider.GetUser(ctx, []interface{}{"A001"}, user)
```

On the provider side, a `context.Context` is passed as the first parameter along with the user data. Here below the code from the current samples shows how the service method is implemented, and how the user data is fetched from the attachment:

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

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.

