# Dubbo-go Config-Center Sample

## 1. Introduction

This example shows dubbo-go's dynamic configuration feature with Apollo as config-center.

## 2. How to run

### Configure the configuration file into apollo

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

Open `https://localhost:8070` with browser, make sure the relevant configuration is already in place in apollo.

### Start an instance with apollo as the configuration center

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
Greet response: greeting:"Hello, Apollo"
```