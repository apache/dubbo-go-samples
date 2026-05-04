# Dubbo-go Config-Center Sample

## 1. 介绍

本示例演示Dubbo-Go以ZooKeeper为配置中心来实现动态配置功能

## 2. 如何运行

### 启动 Zookeeper 实例

确保有一个 Zookeeper 实例监听在 `127.0.0.1:2181`。最简单的方式是用 Docker 启动一个：

```
docker run -d --name zookeeper -p 2181:2181 zookeeper:3.8
```

或者参考 [Zookeeper 安装文档](https://zookeeper.apache.org/doc/current/zookeeperStarted.html) 在本地安装。

### 把配置文件配置到zookeeper中

```yaml
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: '127.0.0.1:2181'
  protocols:
    triple:
      name: tri
      port: 50000
  provider:
    services:
      GreeterProvider:
        interface: greet.GreetService
```

打开本地ZooKeeper客户端查看配置是否成功。
如果没有预先配置也没关系，示例代码中包含了先将配置推送到配置中心的逻辑。

### 以zookeeper作为配置中心启动一个实例

```go
zkOption := config_center.WithZookeeper()
dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-go-server")
addressOption := config_center.WithAddress("127.0.0.1:2181")
groupOption := config_center.WithGroup("dubbogo")
ins, err := dubbo.NewInstance(
    dubbo.WithConfigCenter(zkOption, dataIdOption, addressOption, groupOption),
)
if err != nil {
    panic(err)
}
```

### 启动服务端并注册服务

```go
srv, err := ins.NewServer()
if err != nil {
    panic(err)
}

if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
    panic(err)
}

if err := srv.Serve(); err != nil {
    logger.Error(err)
}
```

### 运行go服务端
```
go run ./go-server/cmd/main.go
```

### 运行go客户端

```
go run ./go-client/cmd/main.go
```

### 在运行Java服务端/客户端之前
```
mvn clean compile
```

### 运行Java服务端（windows）
```
mvn -pl java-server exec:java "-Dexec.mainClass=org.example.server.ZookeeperJavaServer" 
```

### 运行Java客户端（windows）
```
mvn -pl java-client exec:java "-Dexec.mainClass=org.example.client.ZookeeperJavaClient"
```

### 预期的输出
Go/Java 客户端输出：
```
Server response: Hello, this is dubbo go/java server! I received: Hello, this is dubbo go/java client!
```
Go/Java 服务端输出：
```
Received request: Hello, this is dubbo go/java client!
```