# Dubbo-go Config-Center Sample

## 1. 介绍

本示例演示Dubbo-Go以apollo为配置中心来实现动态配置功能

## 2. 如何运行

### 把配置文件配置到apollo中

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

使用浏览器打开`https://localhost:8070` ，确保apollo中已有相关配置。

### 以apollo作为配置中心启动一个实例

```go
ins, err := dubbo.NewInstance(
    dubbo.WithConfigCenter(
        config_center.WithApollo(),
        config_center.WithAddress("127.0.0.1:8080"),
        config_center.WithNamespace("dubbo.yml"),
        config_center.WithDataID("dubbo.yml"),
        config_center.WithAppID("SampleApp"),
        config_center.WithCluster("default"),
        config_center.WithFileExtProperties(),
    ),
)
if err != nil {
    logger.Fatal(err)
}
```

### 启动服务端并注册服务

```go
srv, err := ins.NewServer()
if err != nil {
    logger.Fatal(err)
}

if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
    logger.Fatal(err)
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
Greet response: greeting:"Hello, Apollo"
```