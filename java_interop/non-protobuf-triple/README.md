# Dubbo java and go interoperability, non-protobuf and triple protocol

Please refer to [Multiple Protocols](https://github.com/apache/dubbo-go-samples/tree/main/rpc/multi-protocols) for how to write non-protobuf style protocol.

# Run

## Java

client
```bash
# in directory java-client
mvn -e clean compile exec:java -Dexec.mainClass="org.apache.dubbo.tri.hessian2.client.Application"
```

server
```bash
# in directory java-server
mvn -e clean compile exec:java -Dexec.mainClass="org.apache.dubbo.tri.hessian2.provider.Application"
```

## go

client
```bash
go run go-client/client.go
```

server
```bash
go run go-server/server.go
```
