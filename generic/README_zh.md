# 泛化调用示例

[English](README.md) | [中文](README_zh.md)

本示例演示了如何使用 Dubbo 和 Triple 协议进行泛化调用，实现 Go 和 Java 服务之间的互操作。泛化调用允许在没有服务接口定义的情况下调用远程服务。

## 目录结构

```
generic/
├── go-server/      # Go 服务端（Triple 协议，端口 50052）
├── go-client/      # Go 客户端，泛化调用（直连模式）
├── java-server/    # Java 服务端（Triple 协议，端口 50052）
└── java-client/    # Java 客户端，泛化调用
```

## 前置条件

启动 ZooKeeper（服务端注册服务时需要）：

```bash
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.8
```

## 启动 Go 服务端

```bash
cd generic/go-server/cmd
go run .
```

服务端通过 Triple 协议监听 `50052` 端口，注册到 ZooKeeper，提供 `UserProvider` 服务（version=1.0.0，group=triple）。

## 启动 Go 客户端

```bash
cd generic/go-client/cmd
go run .
```

客户端使用直连模式（`client.WithURL`）连接服务端，通过 `cli.NewGenericService` 进行泛化调用。同时测试 Dubbo 协议（端口 20000）和 Triple 协议（端口 50052）。

## 启动 Java 服务端

在 java-server 目录下构建并运行：

```bash
cd generic/java-server
mvn clean compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.ApiProvider"
```

## 启动 Java 客户端

```bash
cd generic/java-client
mvn clean compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.ApiTripleConsumer"
```

客户端使用 `reference.setGeneric("true")` 进行泛化调用。

## 测试方法

| 方法 | 参数 | 返回值 |
|------|------|--------|
| GetUser1 | String | User |
| GetUser2 | String, String | User |
| GetUser3 | int | User |
| GetUser4 | int, String | User |
| GetOneUser | - | User |
| GetUsers | String[] | User[] |
| GetUsersMap | String[] | Map<String, User> |
| QueryUser | User | User |
| QueryUsers | User[] | User[] |
| QueryAll | - | Map<String, User> |

## 预期输出

服务端日志：

```
Generic Go server started on port 50052
Registry: zookeeper://127.0.0.1:2181
```

客户端日志：

```
[Triple] GetUser1(userId string) res: {id=A003, name=Joe, age=48, ...}
[Triple] GetUser2(userId string, name string) res: {id=A003, name=lily, age=48, ...}
...
All generic call tests completed
```

## 注意事项

- 不要同时启动 Go 服务端和 Java 服务端，它们都监听 50052 端口。
- Go 服务端需要 ZooKeeper 进行服务注册。
- Go 客户端使用直连模式，不依赖 ZooKeeper。
