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
cd http3
go run ./go-server/cmd
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
cd http3
go run ./go-client/cmd
```

### 启动Java服务端

> Java 模块使用 dubbo `3.3.7-SNAPSHOT`，其中包含 HTTP/3 loopback 修复
> （[apache/dubbo#16463](https://github.com/apache/dubbo/pull/16463)）。更早的版本会把显式配置的
> `127.0.0.1` 重写成其它本地地址，而 QUIC 要求客户端拨号的地址在握手期间保持稳定，
> 因此在多网卡机器上 Java 与 Java 之间的握手会超时。`3.3.7` 正式发布后，请改为固定该版本
> 并移除 pom 中的 `apache-snapshots` 仓库。

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

- Go 服务端和客户端通过代码选项 `triple.WithHttp3Enable()` 启用 HTTP/3；默认开启 Alt-Svc 协商，可用 `triple.WithHttp3Negotiation(false)` 关闭。
- Java 服务端和客户端通过 `dubbo.properties` 启用 HTTP/3：
  - `dubbo.protocol.triple.http3.enabled=true` - 启用 HTTP/3 协议
  - `dubbo.protocol.triple.http3.negotiation=false` - 禁用协议协商（强制使用 HTTP/3）
- TLS 证书配置用于安全的 QUIC 连接（见 `x509` 目录）

### QUIC 传输调优

Go 服务端和客户端均通过 `triple.WithHttp3*` 代码选项配置 QUIC 传输参数：

- `WithHttp3KeepAlivePeriod(30 * time.Second)` - 每 30s 发送一次 QUIC keep-alive 包
- `WithHttp3MaxIdleTimeout(90 * time.Second)` - 空闲 QUIC 连接 90s 后关闭
- `WithHttp3MaxIncomingStreams(1024)` / `WithHttp3MaxIncomingUniStreams(1024)` - 并发双向/单向流数量上限
- `WithHttp3InitialStreamReceiveWindow(512 * 1024)` / `WithHttp3MaxStreamReceiveWindow(2 * 1024 * 1024)` - Stream 级流控接收窗口
- `WithHttp3InitialConnectionReceiveWindow(2 * 1024 * 1024)` / `WithHttp3MaxConnectionReceiveWindow(8 * 1024 * 1024)` - Connection 级流控接收窗口

所有选项均可选，未设置时回退到 quic-go 默认值。

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
