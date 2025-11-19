# dubbo-go Helloworld 示例

本示例演示了 dubbo-go 作为 RPC 框架的基本用法,并展示了 Java ↔ Go 之间的跨语言调用。详细说明请参考[官方文档-快速开始](https://dubbo.apache.org/zh-cn/overview/mannual/golang-sdk/quickstart/)。

## 目录结构

- go-server/cmd/main.go - 服务端主程序，包含服务定义、处理器和 RPC 服务端
- go-client/cmd/main.go - RPC 客户端
- java-server/src/main/java/org/example/server/JavaServerApp.java - Java服务端
- java-client/src/main/java/org/example/client/JavaClientApp.java - Java客户端
- proto - API 的 protobuf 定义

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
go run ./go-server/cmd/main.go
```

测试服务端是否正常：
```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

### 启动Golang客户端
```shell
go run ./go-client/cmd/main.go
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
mvn exec:java -Dexec.mainClass=org.example.server.JavaServerApp
```

**在 Windows PowerShell 上：**
```powershell
cd java-server
mvn exec:java "-Dexec.mainClass=org.example.server.JavaServerApp"
```

**或使用提供的脚本（Linux/Mac）：**
```shell
cd java-server
./run.sh
```

测试服务端是否正常：
```shell
curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

### 启动Java客户端

启动 Java 客户端：

**在 Linux/Mac/Git Bash 上：**
```shell
cd java-client
mvn exec:java -Dexec.mainClass=org.example.client.JavaClientApp
```

**在 Windows PowerShell 上：**
```powershell
cd java-client
mvn exec:java "-Dexec.mainClass=org.example.client.JavaClientApp"
```

**或使用提供的脚本（Linux/Mac）：**
```shell
cd java-client
./run.sh
```
## 注意
不能同时启动 Go 和 Java 服务端,Go 服务端 和 Java 服务端 都监听 相同的端口：20000，并暴露 相同的 Triple 服务路径：greet.GreetService/Greet

[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler 安装]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/
[Maven]: https://maven.apache.org/download.cgi
