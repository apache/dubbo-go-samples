# TLS 示例

## 说明

本示例演示了如何在 Dubbo-Go 中使用 TLS（基于 X.509 证书）实现客户端与服务端之间的加密通信和/或双向认证。示例包括一个简单的 `greet` 服务、客户端与服务端的示例程序，以及用于生成测试证书的 X.509 脚本和示例证书。

## 目录结构

* **client/**: 客户端示例程序
* **server/**: 服务端示例程序
* **proto/**: `greet` 服务的 Proto 文件及生成的代码
* **x509/**: 生成和存放测试证书的脚本和示例证书

## 先决条件

* Go 语言 (建议使用 1.18+ 版本)
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

### 1. 启动服务端

在项目根目录下执行以下命令来启动服务端：

```bash
go run ./tls/server/cmd  
```

服务端会加载 `x509/` 目录下的服务器证书和 CA，默认监听配置中指定的地址。如果需要自定义配置，请修改 `server` 程序或源代码中的相关内容。

### 2. 启动客户端

在另一个终端执行以下命令来启动客户端：

```bash
go run ./tls/client/cmd  
```

客户端会使用 `x509/` 目录下的客户端证书和 CA 与服务端建立 TLS 连接，并调用 `greet` 服务。

## 注意事项

* 证书路径和是否启用双向认证的设置取决于示例程序中加载的文件。请查看 `tls/server/cmd/main.go` 和 `tls/client/cmd/main.go` 以了解具体行为和可用的命令行参数。
* 在 Windows 环境下运行 `create.sh` 脚本时，可能需要 WSL/Git Bash，或手动执行 OpenSSL 命令。
* 本示例用于教学和测试目的，示例证书不应在生产环境中使用。
