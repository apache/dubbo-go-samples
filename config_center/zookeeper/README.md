# Dubbo-go Config-Center Sample

## 1. Introduction

This example shows dubbo-go's dynamic configuration feature with Zookeeper as config-center.

## 2. How to run

### Add configuration to the configuration center of zookeeper

```go
dynamicConfig, err := config.NewConfigCenterConfigBuilder().
SetProtocol("zookeeper").
SetAddress("127.0.0.1:2181").
Build().GetDynamicConfiguration()
if err != nil {
    panic(err)
}

if err := dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-server", "dubbogo", configCenterZKServerConfig); err != nil {
    panic(err)
}
time.Sleep(time.Second * 10)
```

Open the local zookeeper client to see if the configuration is successful

### Start an instance with zookeeper as the configuration center

```go
zkOption := config_center.WithZookeeper()
dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-server")
addressOption := config_center.WithAddress("127.0.0.1:2181")
groupOption := config_center.WithGroup("dubbogo")
ins, err := dubbo.NewInstance(
    dubbo.WithConfigCenter(zkOption, dataIdOption, addressOption, groupOption),
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