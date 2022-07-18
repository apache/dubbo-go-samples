# Skywalking with Dubbo-go 3.0

api 定义以及 go 客户端、服务端启动，可以参考 [dubbo-go 3.0 quickstart](https://dubbo.apache.org/zh/docs/quick-start/)

## 使用方法
1. 启动 nacos

2. 启动服务端

使用 goland 启动 skywalking/go-server/cmd/server. 注意需要将`server.go`中的 `YOUR_SKYWALKING_DOMAIN_NAME_OR_IP` 修改为你真实的环境配置，以及修改 skywalking/go-server/conf/`dubbogo.yml`中的内容。

3. 启动客户端

使用 goland 启动 skywalking/go-client/cmd/client. 注意需要将`client.go`中的 `YOUR_SKYWALKING_DOMAIN_NAME_OR_IP` 修改为你真实的环境配置，以及修改 skywalking/go-client/conf/`dubbogo.yml`中的内容。

