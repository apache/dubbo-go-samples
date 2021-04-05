# Dubbo-go context 使用demo 
### 1. 介绍

可以在服务端使用context读取dubbo-go框架默认提供的字段。

也可以使用ctx从客户端传递用户希望传递的字段到服务端。

### 2. 如何获取dubbogo提供的默认字段
你可以获取dubbo attachement 字段：
```go
ctxAtta := ctx.Value(constant.DubboCtxKey("attachment")).(map[string]interface{})
	rsp := ContextContent{
		Path:          ctxAtta["path"].(string),
		InterfaceName: ctxAtta["interface"].(string),
		DubboVersion:  ctxAtta["dubbo"].(string),
		LocalAddr:     ctxAtta["local-addr"].(string),
		RemoteAddr:    ctxAtta["remote-addr"].(string),
	}
```

### 3. 从客户端通过context传递你想要的字段
- 客户端
  如样例 go-client/cliemt.go 所展示：

```go
    // 创建request
    rspContent := &pkg.ContextContent{}

    // 创建 attachment, 类型必须为 map[string]interface{}
	atta := make(map[string]interface{})
    // 添加你想传递的字段
	atta["string-value"] = "string-demo"
	atta["int-value"] = 1231242
    // 需要保证你想传递的结构体提前被注册在了hessian2上
	atta["user-defined-value"] = pkg.ContextContent{InterfaceName: "test.interface.name"}

    // 使用构建好字段的context作为request context
	reqContext := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), atta)
	err := userProvider.GetContext(reqContext, []interface{}{"A001"}, rspContent)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %+v\n", rspContent)
```
- 服务端:
  如样例在 go-server/server.go 所展示

```go
    // 从ctx 获取 attachment
    ctxAtta := ctx.Value(constant.DubboCtxKey("attachment")).(map[string]interface{})
    
    // 从ctx获取用户自定义结构体
	userDefinedval := ctxAtta["user-defined-value"].(*ContextContent)
    
    // 获取特定字段值
	intValue := ctxAtta["int-value"].(int64),
```

请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。


