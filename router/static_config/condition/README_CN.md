# 静态条件路由

这个示例展示了如何在代码中直接配置 Dubbo-Go 的静态条件路由，并通过直连地址访问服务实例。

[English](README.md) | 中文

## 前置准备

- Go 1.25+。

## 这个示例演示了什么

- 两个通过直连地址访问的 provider，分别监听 `20000` 和 `20001`
- 使用 `dubbo.WithRouter(...)` 配置服务级静态 condition router
- 通过 `force=true` 规则把 `Greet` 请求固定路由到 `20000`

## 如何运行

### 启动 Provider

在两个终端中分别启动下面两个 provider：

```shell
$ go run ./go-server/cmd/server.go
$ go run ./go-node2-server/cmd/server_node2.go
```

- `go-server` 监听 `:20000`
- `go-node2-server` 监听 `:20001`

### 启动 Consumer

如果你需要修改 provider 地址，请同步修改 `go-client/cmd/client.go` 中的 `directURL`。

```shell
$ go run ./go-client/cmd/client.go
```

客户端会通过直连 URL 连接两个 provider。
这个示例不需要注册中心和配置中心。

## 预期结果

- 客户端会打印 `invoke successfully: receive: static condition router, response from: server-node-20000`

## 关键路由配置

consumer 注入的是一个服务级静态 condition router：

```go
dubbo.WithRouter(
    router.WithScope("service"),
    router.WithKey(greet.GreetServiceName),
    router.WithForce(true),
    router.WithConditions([]string{
        "method = Greet => port = 20000",
    }),
)
```
