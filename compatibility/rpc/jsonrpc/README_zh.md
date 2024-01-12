# JSON-RPC 示例

## 背景

Dubbo3 提供了 Triple(Dubbo3)、Dubbo2 协议，这是 Dubbo 框架的原生协议。
除此之外，Dubbo3 也对众多第三方协议进行了集成，并将它们纳入 Dubbo 的编程与服务治理体系，
包括 gRPC、Thrift、**JSON-RPC**、Hessian2、REST 等。以下重点介绍 **JSON-RPC** 协议示例。

## 如何启动

- 启动注册中心
- 启动 go-server、go-client 查看 **JSON-RPC** 效果
- 启动 java-server、java-client 查看 **JSON-RPC** 效果

### 启动注册中心

启动项目提供注册中心

```shell
docker-compose -f go-server/docker/docker-compose.yml up -d
```

关闭注册中心

```shell
docker-compose -f go-server/docker/docker-compose.yml dowm
```

### 启动 Go Server、Client

注：Goland 用户可以直接使用 `.run` 配置的启动方式, 详情参考 [HOWTO.md](../HOWTO_zh.md)

启动 go-server：

配置 Dubbogo 配置文件路径（[server/dubbogo.yml](go-server/conf/dubbogo.yml)）：

```shell
DUBBO_GO_CONFIG_PATH=${$PROJECT_DIR$}/dubbo-go-samples/rpc/jsonrpc/go-server/conf/dubbogo.yml
```

启动 go-client：

配置 Dubbogo 配置文件路径（[client/dubbogo.yml](go-client/conf/dubbogo.yml)）：

```shell
DUBBO_GO_CONFIG_PATH=${$PROJECT_DIR$}/dubbo-go-samples/rpc/jsonrpc/go-client/conf/dubbogo.yml
```

### 启动 Java Server、Client

启动 java-server：

可直接运行项目提供 [build.sh](java-server/build.sh) ，基于 maven 环境启动

```shell
bash build.sh
```

启动 java-client：

可直接运行项目提供 [build.sh](java-client/build.sh)，基于 maven 环境启动

```shell
bash build.sh
```


