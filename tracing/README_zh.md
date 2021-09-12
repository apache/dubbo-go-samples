# RPC Dubbo for Dubbo-go 3.0

api 定义以及 go 客户端、服务端启动，可以参考 [dubbo-go 3.0 快速开始](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## 使用方法
1. 启动 `docker/docker-compose.yml` 文件里的 `zookeeper` 服务
2. 启动 server 端：
   1. 找到 go-server 文件夹，运行 cmd 包下的 `main` 函数，启动 golang server
3. 启动 client 端：
   1. 找到 go-client 文件夹，运行 cmd 包下的 `main` 函数，启动 golang client
4. 测试 ziplink步骤：
   1. 启动 `docker/docker-compose.yml` 文件里的 `zipkin` 服务
   2. 启用`go-server`和`go-client` 里的`main` 函数的 `initZipkin()`语句，并注释 `initJaeger()` 语句
   3. 使用 client 调用 server 端服务，用浏览器打开 http://localhost:9411/zipkin 即可看到tracing数据
5. 测试jaeger步骤：
   1. 启动 `docker/docker-compose.yml` 文件里的 `jaeger` 服务
   2. 启用 `go-server` 和 `go-client` 里的 `main` 函数的 `initJaeger()` 语句，并注释 `initZipkin()` 语句
   3. 使用client调用server端服务，用浏览器打开 http://localhost:16686/search 即可看到tracing数据

