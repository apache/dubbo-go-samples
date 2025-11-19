# Helloworld for dubbo-go

This example demonstrates the basic usage of dubbo-go as an RPC framework,and shows Dubbo-Go with Java Interoperability. Check [Quick Start](https://dubbo.apache.org/zh-cn/overview/mannual/golang-sdk/quickstart/) on our official website for detailed explanation.

## Contents

- go-server/cmd/main.go - is the main definition of the service, handler and rpc server
- go-client/cmd/main.go - is the rpc client
- java-server/src/main/java/org/example/server/JavaServerApp.java - is the Java server
- java-client/src/main/java/org/example/client/JavaClientApp.java - is the Java client
- proto - contains the protobuf definition of the API

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

    Generate related stub code with protoc-gen-go and protoc-gen-go-triple:

    ```shell
    protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./proto/greet.proto
    ```
   
4. Install `Maven` [Maven][]

### Run Golang server
```shell
go run ./go-server/cmd/main.go
```

test server work as expected:
```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

### Run Golang client
```shell
go run ./go-client/cmd/main.go
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
mvn exec:java -Dexec.mainClass=org.example.server.JavaServerApp
```

**On Windows PowerShell:**
```powershell
cd java-server
mvn exec:java "-Dexec.mainClass=org.example.server.JavaServerApp"
```

**Or use the provided script (Linux/Mac):**
```shell
cd java-server
./run.sh
```

Test server works as expected:
```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

### Run Java client

Run the Java client:

**On Linux/Mac/Git Bash:**
```shell
cd java-client
mvn exec:java -Dexec.mainClass=org.example.client.JavaClientApp
```

**On Windows PowerShell:**
```powershell
cd java-client
mvn exec:java "-Dexec.mainClass=org.example.client.JavaClientApp"
```

**Or use the provided script (Linux/Mac):**
```shell
cd java-client
./run.sh
```
## Attention
Do NOT Start Go Server and Java Server at the Same Time. Both the Go server and Java server listen on the same port: 20000 and expose the same Triple service path:greet.GreetService/Greet

[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler Installation]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/
[Maven]: https://maven.apache.org/download.cgi
