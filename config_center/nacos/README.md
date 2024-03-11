# Dubbo-go Config-Center Sample

## 1. Introduction

This example shows dubbo-go's dynamic configuration feature with Nacos as config-center.

## 2. How to run

### Add configuration to the configuration center of nacos

```go
dynamicConfig, err := config.NewConfigCenterConfigBuilder().
SetProtocol("nacos").
SetAddress("127.0.0.1:8848").
Build().GetDynamicConfiguration()

if err != nil {
    panic(err)
}

if err := dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-nacos-server", "dubbo", configCenterNacosServerConfig); err != nil {
    panic(err)
}
time.Sleep(time.Second * 10)
```
Open `https://localhost:8848/nacos/` with browser, make sure the relevant configuration is already in place in nacos.

### Start an instance with nacos as the configuration center

```go
nacosOption := config_center.WithNacos()
dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-nacos-server")
addressOption := config_center.WithAddress("127.0.0.1:8848")
groupOption := config_center.WithGroup("dubbo")
ins, err := dubbo.NewInstance(
    dubbo.WithConfigCenter(nacosOption, dataIdOption, addressOption, groupOption),
)
if err != nil {
    panic(err)
}
```

### Start server and register for the service

```go
srv, err := ins.NewServer()
if err != nil {
    panic(err)
}

if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
    panic(err)
}

if err := srv.Serve(); err != nil {
    logger.Error(err)
}
```

### Run client
```shell
$ go run ./go-client/cmd/main.go
```

### Expect output
```
Greet response: greeting:"hello world"
```