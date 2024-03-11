# Dubbo-go Config-Center Sample

## 1. 介绍

This example shows dubbo-go's dynamic configuration feature with Zookeeper as config-center.

## 2. 如何运行

### 向zookeeper中添加配置

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

### 以zookeeper作为配置中心启动一个实例

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

### 启动服务端并注册服务

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

### 启动客户端
```shell
$ go run ./go-client/cmd/main.go
```

### 预期的输出
```
Greet response: greeting:"hello world"
```