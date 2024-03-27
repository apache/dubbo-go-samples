# Zookeeper as registry

This example shows dubbo-go's service discovery and java-go interoperation feature  with Nacos as registry.


> before run the code , you should Follow this instruction to install and start Nacos server.

## interface level discovery
### how to run
#### java server <-> go client 
```shell
cd interface
```
**start java server**
```shell
cd java-server
sh run.sh
```

**start go client**
```shell
cd go-client
go run client.go

```

#### go server <-> java client 
**start go server**
```shell
cd go-server
go run server.go
```
**start java client**
```shell
cd java-client
sh run.sh
```


## service level discovery

```shell
cd service
```
**start java server**
```shell
cd java-server
sh run.sh
```

**start go client**
```shell
cd go-client
go run client.go

```

#### go server <-> java client 
**start go server**
```shell
cd go-server
go run server.go
```
**start java client**
```shell
cd java-client
sh run.sh
```