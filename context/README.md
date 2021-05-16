# Dubbo-go context Usage 
### 1. Introduction

Context in server end can be used to read specific field that dubbo-go framework provided by default.\
It can also used to transfer base-type of golang, even 

### 2. How to get dubbo-go default field
You can get dubbo attachment in this way:
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

### 3. Transfer value you want from client to server
Client end\
As demo in go-client/cliemt.go shows
```go
    // create requset context
    rspContent := &pkg.ContextContent{}

    // create attachment, which must be map[string]interface{}
	atta := make(map[string]interface{})
    // add fields you like
	atta["string-value"] = "string-demo"
	atta["int-value"] = 1231242
    // make sure the UserDefined Pkg is registered to hessian2
	atta["user-defined-value"] = pkg.ContextContent{InterfaceName: "test.interface.name"}

    // invoke with your context
	reqContext := context.WithValue(context.Background(), constant.DubboCtxKey("attachment"), atta)
	err := userProvider.GetContext(reqContext, []interface{}{"A001"}, rspContent)
	if err != nil {
		gxlog.CError("error: %v\n", err)
		os.Exit(1)
		return
	}
	gxlog.CInfo("response result: %+v\n", rspContent)
```
Server end:\
As demo in go-server/server.go shows
```go
    // get attachment from ctx
    ctxAtta := ctx.Value(constant.DubboCtxKey("attachment")).(map[string]interface{})
    
    // get user defined struct from attachment
	userDefinedval := ctxAtta["user-defined-value"].(*ContextContent)
    
    // get value you sent
	intValue := ctxAtta["int-value"].(int64),
```

Pls. refer to [HOWTO.md](../HOWTO.md) under the root directory to run this sample.