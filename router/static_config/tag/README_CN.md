# 静态标签路由

这个示例展示了如何在代码中直接配置 Dubbo-Go 的静态标签路由，不依赖注册中心和配置中心。

[English](README.md) | 中文

## 前置准备

- Go 1.25+。

## 这个示例演示了什么

- 一个不带 tag 的 provider，监听 `20000`
- 一个带 `gray` tag 的 provider，监听 `20002`
- 使用 `dubbo.WithRouter(...)` 配置应用级静态 tag router

## 如何运行

### 启动 Provider

在两个终端中分别启动下面两个 provider：

```shell
$ go run ./go-server/cmd/server.go
$ go run ./go-tag-server/cmd/server_tag.go
```

- `go-server` 监听 `:20000`，不带 tag
- `go-tag-server` 监听 `:20002`，并以 `gray` 标签导出服务

### 启动 Consumer

```shell
$ go run ./go-client/cmd/client.go
```

客户端只使用直连 URL：

```text
tri://127.0.0.1:20000;tri://127.0.0.1:20002?dubbo.tag=gray
```

这个示例不需要注册中心，也不需要配置中心。

## 预期结果

客户端会运行一个场景，并路由到 `server-with-gray-tag`。
你会看到类似下面的日志：

```text
INFO ... invoke successfully: receive: static tag router, response from: server-with-gray-tag
```

## 关键路由配置

静态 tag router：

```go
dubbo.WithRouter(
    router.WithScope("application"),
    router.WithKey("static-tag-provider"),
    router.WithForce(false),
    router.WithTags([]global.Tag{
        {
            Name:      "gray",
            Addresses: []string{"127.0.0.1:20002"},
        },
    }),
)
```

携带请求 tag：

```go
ctx := context.WithValue(context.Background(), constant.AttachmentKey, map[string]string{
    constant.Tagkey: "gray",
})
```
