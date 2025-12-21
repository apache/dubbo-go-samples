# Dubbo-go Config-Center Sample

## 1. 介绍

本示例演示Dubbo-Go以nacos为配置中心来实现动态配置功能

## 2. 如何运行
### 在 `integrate_test/dockercompose`目录下执行：
```
docker compose up -d nacos zookeeper
```
### 把配置文件配置到nacos中

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
      port: 20000
  provider:
    services:
      GreeterProvider:
        interface: greet.GreetService
```

使用浏览器打开`https://localhost:8848/nacos/` ，确保nacos中已有相关配置。
如果没有预先配置也没关系，示例代码中包含了先将配置推送到配置中心的逻辑。
### 以nacos作为配置中心启动一个实例

```go
nacosOption := config_center.WithNacos()
dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-nacos-server")
addressOption := config_center.WithAddress("127.0.0.1:8848")
groupOption := config_center.WithGroup("dubbo")
ins, err := dubbo.NewInstance(
    dubbo.WithConfigCenter(nacosOption, dataIdOption, addressOption, groupOption),
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
mvn -pl java-server exec:java "-Dexec.mainClass=org.example.server.NacosJavaServer" 
```

### 运行Java客户端（windows）
```
 mvn -pl java-client exec:java "-Dexec.mainClass=org.example.client.NacosJavaClient"
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