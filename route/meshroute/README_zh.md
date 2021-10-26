# Helloworld for Dubbo-go 3.0

api 定义以及 go 客户端、服务端启动，可以参考 [dubbo-go 3.0 快速开始](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## 使用方法
1. 可以使用 docker 结合 `integrate_test/dockercompose/docker-compose.yml` 文件或下载官网的二进制文件启动 zookeeper 服务。
2. 启动服务端服务。
3. 获取动态配置中心，并发布 mesh 路由配置。
4. 使用客户端调用服务端发布的服务。