# Condition router

这个例子展示了如何使用dubbo-go的condition router功能。

[English](README.md) | 中文

## 前置准备

- Docker以及Docker Compose环境来运行Nacos注册中心。
- Go 1.23+版本。
- Nacos 2.x+版本。

## 如何运行

### 启动Nacos注册中心

参考这个教程来[启动Nacos](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/)。

### 运行服务端(Provider)

在这个示例中，你将运行两个服务端，分别在20000以及20001端口上提供服务。

```shell
$ go run ./go-server/cmd/server.go              # 20000端口
$ go run ./go-node2-server/cmd/server_node2.go  # 20001端口
```

### 运行客户端（Consumer）

在这个示例中，客户端将在一个死循环中一直调用Greet方法，你需要：

- 启动客户端，观察其调用时的负载均衡（Load Balance）。
- 在Nacos注册中心上设置`condition router`的配置，再次观察客户端的调用情况。

```shell
$ go run ./go-client/cmd/client.go
```

### Nacos配置

新建一个`Data ID`为`condition-server.condition-router`，格式为`yaml`的配置。

Group设置为`DEFAULT_GROUP`。

> 注意：Nacos中命名规则为{application.name}.{router_type}

```yaml
configVersion: V3.3.2
scope: "application"
key: "condition-server"
priority: 1
force: true
enabled: true
conditions:
  - from:
      match: "application = condition-client"
    to:
      - match: "port = 20001"
```

## 预期结果

- 启动客户端但是未在nacos配置中心设置condition router的配置的时候，客户端将在两个服务端之间来回调用。
- 启动客户端并在nacos配置中心设置了condition router后，客户端将只调用其中一个服务端。



