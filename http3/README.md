# HTTP/3 for dubbo-go

English | [中文](README_CN.md)

This example demonstrates how to use dubbo-go with HTTP/3 protocol support via the Triple protocol. It shows how to enable HTTP/3 for high-performance communication between Go and Java services using TLS for secure connections.

## Contents

- go-server/cmd/main.go - is the main definition of the service, handler and rpc server with HTTP/3 support
- go-client/cmd/main.go - is the rpc client with HTTP/3 support
- java-server/src/main/java/org/apache/dubbo/samples/http3/Http3ServerApp.java - is the Java server with HTTP/3 support
- java-client/src/main/java/org/apache/dubbo/samples/http3/Http3ClientApp.java - is the Java client with HTTP/3 support
- proto - contains the protobuf definition of the API
- x509 - contains TLS certificates and keys for secure connections

## Key Features

- **HTTP/3 Protocol Support**: Uses QUIC transport for faster, more reliable connections
- **Cross-Language Interoperability**: Demonstrates Go and Java interoperability with HTTP/3
- **TLS Encryption**: Secure communication with client and server certificates
- **Triple Protocol**: Built on Apache Dubbo's Triple protocol with HTTP/3 enabled

## How to run

### Prerequisites
1. Install `protoc` [version3][]
   Please refer to [Protocol Buffer Compiler Installation][].

2. Install `protoc-gen-go` and `protoc-gen-triple`
   Install the version of your choice of protoc-gen-go. here use the latest version as example:

    ```shell
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
    ```
   
    Install the latest version of protoc-gen-triple:

    ```shell
    go install github.com/dubbogo/protoc-gen-go-triple/v3@v3.0.2
    ```

3. Generate stub code

    Generate related stub code with protoc-gen-go and protoc-gen-triple:

    ```shell
    protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./proto/greet.proto
    ```
   
4. Install `Maven` [Maven][]

### Run Golang server
```shell
cd go-server/cmd
go run main.go
```

Test server works as expected:
```shell
curl -k \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    https://localhost:20000/greet.GreetService/Greet
```

### Run Golang client
```shell
cd go-client/cmd
go run main.go
```

### Run Java server

Build all Java modules from the root directory:
```shell
mvn clean compile
```

Run the Java server:

**On Linux/Mac/Git Bash:**
```shell
cd java-server
mvn exec:java -Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ServerApp
```

**On Windows PowerShell:**
```powershell
cd java-server
mvn exec:java "-Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ServerApp"
```

**Or use the provided script (Linux/Mac):**
```shell
cd java-server
./run.sh
```

Test server works as expected:
```shell
curl -k \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    https://localhost:20000/greet.GreetService/Greet
```

### Run Java client

Run the Java client:

**On Linux/Mac/Git Bash:**
```shell
cd java-client
mvn exec:java -Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ClientApp
```

**On Windows PowerShell:**
```powershell
cd java-client
mvn exec:java "-Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ClientApp"
```

**Or use the provided script (Linux/Mac):**
```shell
cd java-client
./run.sh
```

## Configuration

### HTTP/3 Enabled Configuration

The services are configured with HTTP/3 support. Key configuration parameters:

- `protocol.triple.http3.enabled=true` - Enables HTTP/3 protocol
- `protocol.triple.http3.negotiation=false` - Disables protocol negotiation (forces HTTP/3)
- TLS certificates are configured for secure QUIC connections

### Certificate Files

The x509 directory contains the following certificate files:
- `server2_cert.pem` - Server certificate
- `server2_key_pkcs8.pem` - Server private key (PKCS8 format)
- `server_ca_cert.pem` - CA certificate for verification

## Attention

Do NOT start Go Server and Java Server at the same time. Both the Go server and Java server listen on the same port: 20000 and expose the same Triple service path: greet.GreetService/Greet

[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler Installation]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/
[Maven]: https://maven.apache.org/download.cgi
