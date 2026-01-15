# 泛化调用示例 (Triple 协议)

[English](README.md) | [中文](README_zh.md)

本示例演示了如何使用 Triple 协议进行泛化调用，实现 Go 和 Java 服务之间的互操作。泛化调用允许在没有服务接口定义的情况下调用远程服务。

## 目录结构

```
generic/
├── go-server/      # Go 服务端，监听 :50052
├── go-client/      # Go 客户端，泛化调用
├── java-server/    # Java 服务端，监听 :50052
└── java-client/    # Java 客户端，泛化调用
```

## 前置条件

启动 ZooKeeper：

```bash
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.8
```

## 启动 Go 服务端

```bash
cd generic/go-server/cmd
go run .
```

服务端监听 `50052` 端口，并注册到 ZooKeeper。

## 启动 Go 客户端

```bash
cd generic/go-client/cmd
go run .
```

客户端通过 ZooKeeper 发现服务，使用 `client.WithGenericType("true")` 进行泛化调用。

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
Generic Go/Java server started on port 50052
Registry: zookeeper://127.0.0.1:2181
```

客户端日志：

```
[PASS] GetUser1(String): {id=A003, name=Joe, age=48, ...}
[PASS] GetUser2(String, String): {id=A003, name=lily, age=48, ...}
...
[OK] All tests passed!
```

## 注意事项

- 不要同时启动 Go 服务端和 Java 服务端，它们都监听 50052 端口。
- 启动服务端或客户端之前，请确保 ZooKeeper 正在运行。
