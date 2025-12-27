# Dubbo-go Config-Center Sample

## 1. Introduction

This example shows dubbo-go's dynamic configuration feature with Zookeeper as config-center.

## 2. How to run
### Run the following commands under `integrate_test/dockercompose`:

```
docker compose up -d zookeeper
```
### Configure the configuration file into zookeeper

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
        interface: greet.GreetService
```

Open the local zookeeper client to see if the configuration is successful.
If there is no preexisting configuration, that is fine as well, because the sample code already includes logic to first push the configuration to the config center.

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

### Run go server

```
$ go run ./go-server/cmd/main.go
```

### Run go client

```
$ go run ./go-client/cmd/main.go
```
### Before run java server/client
```
mvn clean compile
```
### Run Java server(Windows)
```
mvn -pl java-server exec:java "-Dexec.mainClass=org.example.server.ZookeeperJavaServer" 
```

### Run Java client(Windows)
```
mvn -pl java-client exec:java "-Dexec.mainClass=org.example.client.ZookeeperJavaClient"
```


### Expect output

go/java client output:
```
Server response: Hello, this is dubbo go/java server! I received: Hello, this is dubbo go/java client!
```
go/java server output:
```
Received request: Hello, this is dubbo go/java client!
```