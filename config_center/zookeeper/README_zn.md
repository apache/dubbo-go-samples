# Dubbo-go Config-Center Sample

## 1. 介绍

本示例演示Dubbo-Go以ZooKeeper为配置中心来实现动态配置功能

## 2. 如何运行

### 把配置文件配置到zookeeper中

```yaml
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: '127.0.0.1:2181'
  protocols:
    triple:
      name: tri
      port: 20000
  provider:
    services:
      GreeterProvider:
        interface: com.apache.dubbo.sample.basic.IGreeter
```

打开本地ZooKeeper客户端查看配置是否成功

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