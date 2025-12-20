# HTTP/3 for dubbo-go

[English](README.md) | 中文

本示例演示了如何通过 Triple 协议在 dubbo-go 中使用 HTTP/3 协议支持。展示了使用 HTTP/3 为 Go 和 Java 服务之间的高性能通信提供支持，并使用 TLS 进行安全连接。

## 目录结构

- go-server/cmd/main.go - 服务端主程序，包含服务定义、处理器和 RPC 服务端，支持 HTTP/3
- go-client/cmd/main.go - RPC 客户端，支持 HTTP/3
- java-server/src/main/java/org/apache/dubbo/samples/http3/Http3ServerApp.java - Java服务端，支持 HTTP/3
- java-client/src/main/java/org/apache/dubbo/samples/http3/Http3ClientApp.java - Java客户端，支持 HTTP/3
- proto - API 的 protobuf 定义
- x509 - 安全连接的 TLS 证书和密钥

## 主要特性

- **HTTP/3 协议支持**：通过 QUIC 传输实现更快、更可靠的连接
- **跨语言互操作性**：演示 Go 和 Java 通过 HTTP/3 的互操作性
- **TLS 加密**：使用客户端和服务器证书进行安全通信
- **Triple 协议**：基于 Apache Dubbo 的 Triple 协议，启用 HTTP/3 支持

## 运行方法

### 前置条件
1. 安装 `protoc` [version3][]
   参考[Protocol Buffer Compiler 安装][]。

2. 安装 `protoc-gen-go` 和 `protoc-gen-triple`
   以最新版本为例：

    ```shell
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
    ```
   
    安装 protoc-gen-triple：

    ```shell
    go install github.com/dubbogo/protoc-gen-go-triple/v3@v3.0.2
    ```

3. 生成代码存根

    使用 protoc-gen-go 和 protoc-gen-go-triple 生成相关代码：

    ```shell
    protoc --go_out=. --go_opt=paths=source_relative --go-triple_out=. --go-triple_opt=paths=source_relative ./proto/greet.proto
    ```
   
4. 安装 `Maven` [Maven][]

### 启动Golang服务端
```shell
cd go-server/cmd
go run main.go
```

测试服务端是否正常：
```shell
curl -k \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    https://localhost:20000/greet.GreetService/Greet
```

### 启动Golang客户端
```shell
cd go-client/cmd
go run main.go
```

### 启动Java服务端

从根目录构建所有 Java 模块：
```shell
mvn clean compile
```

启动 Java 服务端：

**在 Linux/Mac/Git Bash 上：**
```shell
cd java-server
mvn exec:java -Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ServerApp
```

**在 Windows PowerShell 上：**
```powershell
cd java-server
mvn exec:java "-Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ServerApp"
```

**或使用提供的脚本（Linux/Mac）：**
```shell
cd java-client
./run.sh
```

测试服务端是否正常：
```shell
curl -k \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    https://localhost:20000/greet.GreetService/Greet
```

### 启动Java客户端

启动 Java 客户端：

**在 Linux/Mac/Git Bash 上：**
```shell
cd java-client
mvn exec:java -Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ClientApp
```

**在 Windows PowerShell 上：**
```powershell
cd java-client
mvn exec:java "-Dexec.mainClass=org.apache.dubbo.samples.http3.Http3ClientApp"
```

**或使用提供的脚本（Linux/Mac）：**
```shell
cd java-client
./run.sh
```

## 配置说明

### HTTP/3 启用配置

服务配置了 HTTP/3 支持。主要配置参数：

- `protocol.triple.http3.enabled=true` - 启用 HTTP/3 协议
- `protocol.triple.http3.negotiation=false` - 禁用协议协商（强制使用 HTTP/3）
- TLS 证书配置用于安全的 QUIC 连接

### 证书文件

x509 目录包含以下证书文件：
- `server2_cert.pem` - 服务器证书
- `server2_key_pkcs8.pem` - 服务器私钥（PKCS8 格式）
- `server_ca_cert.pem` - CA 证书，用于验证

## 注意

不能同时启动 Go 和 Java 服务端。Go 服务端和 Java 服务端都监听相同的端口：20000，并暴露相同的 Triple 服务路径：greet.GreetService/Greet

[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler 安装]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/
[Maven]: https://maven.apache.org/download.cgi
