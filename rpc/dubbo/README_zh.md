# RPC Dubbo for Dubbo-go 3.0

api 定义以及 go 客户端、服务端启动，可以参考 [dubbo-go 3.0 快速开始](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## 使用方法
1. 启动 zookeeper

2. 启动服务端

使用 goland 启动 rpc-dubbo-go-server

或者

在 java-server 文件夹下执行 `rpc/dubbo/java-server/2.6/src/main/java/org/apache/dubbo/Provider#main()` 启动 java server

3. 启动客户端

使用 goland 启动 rpc-dubbo-go-client

或者

在 java-client 文件夹下执行 `rpc/dubbo/java-client/2.6/src/main/java/org/apache/dubbo/Consumer#main()` 启动 java client
