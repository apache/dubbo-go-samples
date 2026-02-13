# Tag router
这个例子展示了如何使用dubbo-go的tag router功能。

## 前置准备

- Docker以及Docker Compose环境来运行Nacos注册中心。
- Go 1.23+版本。

## 如何运行

### 启动Nacos注册中心
参考这个教程来[启动Nacos](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/)。

### 运行服务端（Provider）
在这个示例中，你将运行一个具有tag的server以及一个不具有tag的server。

```shell
$ go run ./go-server/cmd/server-base/server.go      # 无标签server
$ go run ./go-server/cmd/server-tag/server_tag.go   # 有标签server
```

### 运行客户端（Consumer）

```shell
$ go run ./go-client/cmd/client.go
```

## 预期结果

- ✔ invoke successfully : receive: tag with force, response from: server-with-tag
- ❌ invoke failed: Failed to invoke the method Greet.
- ✔ invoke successfully : receive: tag with no-force, response from: server-without-tag
- ✔ invoke successfully : receive: non-tag, response from: server-without-tag
