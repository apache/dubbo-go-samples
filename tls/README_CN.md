# TLS 示例

## 说明

本示例演示了如何在 Dubbo-Go 和 Dubbo-Java 中使用 TLS（基于 X.509 证书）实现客户端与服务端之间的加密通信和/或双向认证。更重要的是，本示例展示了 **Dubbo-Go 与 Dubbo-Java 的跨语言互通能力**— Go 客户端可以无缝调用 Java 服务端，Java 客户端也可以调用 Go 服务端。示例包括 Go 和 Java 的客户端与服务端的示例程序，以及用于生成测试证书的 X.509 脚本和示例证书。

## 目录结构

* **client/**: Go 客户端示例程序
* **server/**: Go 服务端示例程序
* **java-server/**: Java 服务端示例程序（Dubbo-Java 提供者）
* **java-client/**: Java 客户端示例程序（Dubbo-Java 消费者）
* **proto/**: `greet` 服务的 Proto 文件及生成的代码
* **x509/**: 生成和存放测试证书的脚本和示例证书

## 先决条件

* Go 语言 (建议使用 1.18+ 版本)，用于运行 Go 客户端和服务端
* Java 开发环境 (JDK 8+) 和 Maven (3.6.0+)，仅当运行 Java 服务端或客户端时需要
* 在 Windows 系统中，建议使用 Git Bash 或 WSL 来运行证书生成脚本（`x509/create.sh` 使用 OpenSSL）。

## 生成测试证书

1. 进入 `x509` 目录并运行证书生成脚本：

    * 在类 Unix 环境中，执行以下命令：

      ```bash
      cd tls/x509 && ./create.sh  
      ```
    * 在 Windows 系统中，如果未安装 Bash 或 OpenSSL，可以通过 WSL/Git Bash 来运行，或者使用 OpenSSL 手动按照 `x509/openssl.cnf` 文件生成证书。
2. 生成的证书将保存在 `x509/` 目录下，文件包括：

    * `server_ca_*.pem`
    * `client_ca_*.pem`
    * `server{1,2}_*.pem`
    * `client{1,2}_*.pem`

## 如何运行

### 方式一：Go 服务端 + Go 客户端

#### 1. 启动 Go 服务端

在 tls 目录下执行以下命令来启动 Go 服务端：

```bash
go run ./server/cmd  
```

服务端会加载 `x509/` 目录下的服务器证书和 CA，默认监听配置中指定的地址。如果需要自定义配置，请修改 `server` 程序或源代码中的相关内容。

#### 2. 启动 Go 客户端

在另一个终端于 tls 目录下执行以下命令来启动 Go 客户端：

```bash
go run ./client/cmd  
```

客户端会使用 `x509/` 目录下的证书与服务端建立 TLS 连接，并调用 `greet` 服务。

### 方式二：Java 服务端 + Go 客户端（Dubbo-Go ↔ Dubbo-Java 互通）

本方式展示了 **Dubbo-Go 与 Dubbo-Java 通过 Triple 协议和 TLS 进行互通**。

#### 1. 启动 Java 服务端（Dubbo-Java 提供者）

在 tls 目录下，进入 `java-server` 子目录并启动 Maven 项目：

```bash
cd ./java-server
mvn clean compile
mvn exec:java -Dexec.mainClass="org.apache.dubbo.samples.tls.provider.TlsTriProvider"
```

Java 服务端将在 20000 端口启动，使用 `x509/server2_cert.pem` 和 `x509/server2_key.pem` 证书进行 TLS 加密。

#### 2. 启动 Go 客户端

在另一个终端，在 tls 目录下运行 Go 客户端：

```bash
go run ./client/cmd
```

Go 客户端将通过 TLS 连接到 Java 服务端并调用 `greet` 服务。你应该会看到类似如下的输出：

```
Greet response: Hello hello world from Java provider
```

### 方式三：Go 服务端 + Java 客户端

本方式展示了反向的互通——Java 客户端调用 Go 服务端。

#### 1. 启动 Go 服务端

在 tls 目录下执行以下命令来启动 Go 服务端：

```bash
go run ./server/cmd
```

Go 服务端将在 20000 端口启动，使用 TLS 证书进行加密。

#### 2. 启动 Java 客户端

在另一个终端，在 tls 目录下进入 `java-client` 子目录并启动 Maven 项目：

```bash
cd ./java-client
mvn clean compile
mvn exec:java -Dexec.mainClass="org.apache.dubbo.samples.tls.consumer.TlsTriProtoConsumer"
# 如需自定义连接地址与 SNI（证书主机名），可加：-Dtls.host=127.0.0.1 -Dtls.authority=dubbogo.test.example.com
```

Java 客户端将通过 TLS 连接到 Go 服务端并调用 `greet` 服务。你应该会看到类似如下的输出：

```
Greet response: hello world
```

## Dubbo-Go 与 Dubbo-Java 互通性

本示例展示了 **Dubbo-Go 和 Dubbo-Java 框架之间的跨语言互通**：

* **协议**：双方都使用 Triple 协议（基于 gRPC/HTTP2）
* **序列化**：Protobuf 确保语言无关的数据交换
* **TLS/SSL**：双方都支持基于 X.509 证书的加密和认证
* **服务接口**：在 `proto/greet.proto` 中定义，为 Go 和 Java 分别生成代码

## 注意事项

* 证书路径和是否启用双向认证的设置取决于示例程序中加载的文件。请查看 `tls/server/cmd/main.go` 和 `tls/client/cmd/main.go` 以了解具体行为和可用的命令行参数。
* 在 Windows 环境下运行 `create.sh` 脚本时，可能需要 WSL/Git Bash，或手动执行 OpenSSL 命令。
* 本示例用于教学和测试目的，示例证书不应在生产环境中使用。
