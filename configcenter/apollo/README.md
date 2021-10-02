# Apollo config center example


## Instructions


### go-server start

1. Create new apollo server namespace for yaml format

2. Add provider config content 
```yaml
dubbo:
  application:
     name: "demo-server"
     version: "2.0"
  registries:
    "demoZK":
      protocol: "zookeeper"
      timeout: "3s"
      address: "127.0.0.1:2181"
  protocols:
    "triple":
      name: "tri"
      port: 20000
  provider:
    registries:
      - demoZK
    services:
      "greeterImpl":
        protocol: "triple"
        interface: "com.apache.dubbo.sample.basic.IGreeter" # must be compatible with grpc or dubbo-java
```

3. Start provider server

### go-client start

1. Create new apollo client namespace for yaml format

2. Add client config content

```yaml
dubbo:
  registries:
    "demoZK":
      protocol: "zookeeper"
      timeout: "3s"
      address: "127.0.0.1:2181"
  consumer:
    registries:
      - demoZK
    references:
      "greeterImpl":
        protocol: "tri"
        interface: "com.apache.dubbo.sample.basic.IGreeter" # must be compatible with grpc or dubbo-java
```

3. Start provider server



 

