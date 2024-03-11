# Dubbo-go Config-Center Sample

## 1. 介绍

This example shows dubbo-go's dynamic configuration feature with Nacos as config-center.

## 2. 如何运行

### 把配置文件配置到nacos中

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

Open `https://localhost:8848/nacos/` with browser, make sure the relevant configuration is already in place in nacos.

### 以nacos作为配置中心启动一个实例

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