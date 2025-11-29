# 泛化调用

泛化调用是在客户端没有接口信息时保证信息被正确传递的手段，即把 POJO 泛化为通用格式（如字典、字符串），一般被用于集成测试、网关等场景。

本示例演示了 Dubbo-Go 和 Dubbo Java 服务之间的泛化调用，展示了不同语言实现的服务如何互操作。

## 目录结构

- go-server: Dubbo-Go 服务端示例
- go-client: Dubbo-Go 客户端示例（泛化调用）
- java-client: Dubbo Java 客户端示例
- java-server: Dubbo Java 服务端示例
- build: 集成测试需要的脚本

Dubbo Java 示例可以用来测试与 Dubbo-Go 的互操作性。您可以启动 java 服务端配合 go 客户端，或者启动 go 服务端配合 java 客户端进行测试。

## 环境准备

- Docker 和 Docker Compose 用于运行 ZooKeeper 注册中心
- Go 1.23+ 用于 Dubbo-Go 示例
- Java 8+ 和 Maven 用于 Dubbo Java 示例

## 注册中心

本示例使用 ZooKeeper 作为注册中心。以下命令通过 docker 启动 ZooKeeper，因此需要确保已安装 docker 和 docker-compose。

```shell
# 启动 ZooKeeper 注册中心
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.4.14
```

## 运行示例

### Dubbo-Go 服务端

使用 Dubbo-Go 作为服务提供者，可以通过命令行工具启动：

```shell
cd go-server/cmd && go run server.go
```

### Dubbo-Go 客户端（泛化调用）

使用 Dubbo-Go 作为服务消费者进行泛化调用：

```shell
cd go-client/cmd && go run client.go
```

### Dubbo Java 服务端

使用 Dubbo Java 作为服务提供者：

```shell
cd java-server/java-server
mvn clean package
sh run.sh
```

### Dubbo Java 客户端

使用 Dubbo Java 作为服务消费者：

```shell
cd java-client/java-client
mvn clean package
sh run.sh
```

## 测试互操作性

本示例旨在测试 Dubbo-Go 和 Dubbo Java 之间的互操作性：

1. 启动 ZooKeeper 注册中心
2. 启动 go-server 或 java-server 之一
3. 运行 go-client 或 java-client 之一来测试泛化调用

客户端将向服务端发起多种泛化调用，包括：
- GetUser1(String userId)
- GetUser2(String userId, String name)
- GetUser3(int userCode)
- GetUser4(int userCode, String name)
- GetOneUser()
- GetUsers(String[] userIdList)
- GetUsersMap(String[] userIdList)
- QueryUser(User user)
- QueryUsers(List<User> userObjectList)
- QueryAll()