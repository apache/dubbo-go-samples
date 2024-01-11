# gRPC 示例

Dubbo 3.0 提供了 Triple(Dubbo3)、Dubbo2 协议，这两个是 Dubbo 框架的原生协议。 除此之外，Dubbo3 也对众多第三方协议进行了集成，并将它们纳入 Dubbo 的编程与服务治理体系， 包括 gRPC、Thrift、JSON-RPC、Hessian2、REST 等。

**本示例将介绍 gRPC 协议的使用方法**。

## 运行示例：

启动zk，监听127.0.0.1:2181端口。
若未安装 zk , 也可以借助 docker，直接执行下面的命令来启动所有运行samples的依赖组件：zk(2181), nacos(8848), etcd(2379)。

`docker-compose -f {PATH_TO_SAMPLES_PROJECT}/integrate_test/dockercompose/docker-compose.yml up -d`

### 通过命令行运行

- 服务端

`cd rpc/grpc/go-server/cmd` # 进入仓库目录

`export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"`# 设置配置文件环境变量

`go run .` # 启动服务

- 客户端

`cd rpc/grpc/go-client/cmd` # 进入仓库目录

`export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"`# 设置配置文件环境变量

`go run .` # 启动客户端发起调用

### 调用成功

客户端调用后，可看到客户端打印如下信息，调用成功：

`[XXXX-XX-XX/XX:XX:XX main.main: client.go: 55] client response result: {this is message from reply {} [] 0}`
