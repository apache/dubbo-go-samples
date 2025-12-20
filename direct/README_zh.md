# Direct 示例（Triple 直连）

[English](README.md) | [中文](README_zh.md)

在本示例中，Dubbo-Go v3 Triple 服务端直接监听本地端口，客户端通过 `client.WithClientURL("tri://127.0.0.1:20000")` 指定目标地址，完全绕过注册中心，展示最小可运行的点对点调用链路，方便本地调试。

## 目录结构

```
direct/
├── proto/          # greet.proto 以及对应的 triple 代码
├── go-server/      # 提供 greet.GreetService 的Golang服务端
├── go-client/      # 直接拨号 tri://127.0.0.1:20000 的Golang客户端
├── java-server/    # 提供 greet.GreetService 的Java服务端
└── java-client/    # 直接拨号 tri://127.0.0.1:20000 的Java客户端 
```

## 启动Golang服务端

```bash
cd direct/go-server/cmd
go run .
```

服务端监听 `20000` 端口，并实现 `greet.GreetService` 的 `Greet` 方法。

## 启动Golang客户端

```bash
cd direct/go-client/cmd
go run .
```

客户端通过 `client.WithClientURL` 配置直连地址，无需任何 yaml 配置，也无需注册中心即可完成调用。

## 启动Java服务端

从根目录构建所有 Java 模块：
(如没有安装 `Maven` ,请先安装[Maven][])
```shell
mvn clean compile
cd java-server
mvn exec:java "-Dexec.mainClass=org.example.server.JavaServerApp"
```

## 启动Java客户端

从根目录构建所有 Java 模块：
(如没有安装 `Maven` ,请先安装[Maven][])
```shell
mvn clean compile
cd java-client
mvn exec:java "-Dexec.mainClass=org.example.client.JavaClientApp"
```

## 预期输出

服务端日志：

```
INFO ... Direct server form Golang/Java received name = Golang Client dubbo-go/Java Client dubbo
```

客户端日志：

```
INFO ... direct call response: Hello form Java/Golang Server, Golang Client dubbo-go/Java Client dubbo
```

## 注意

由于Golang/Java Server都在监听20000端口,你不能同时打开Golang Server和Java Server




[Maven]: https://maven.apache.org/download.cgi