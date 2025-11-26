# dubbo-go 游戏示例

[English](README.md) | 中文

本示例演示了一个使用 dubbo-go 作为 RPC 框架的足球游戏应用。它包含两个服务：处理游戏逻辑的游戏服务（登录、得分、排名）和提供 Web 前端 HTTP 接口的网关服务，作为前端和游戏服务之间的网关。

## 架构

- **game**: 处理游戏逻辑的游戏服务（Login、Score、Rank）
- **gate**: 为 Web 前端提供 HTTP API 并为游戏服务提供 RPC 服务的网关服务
- **website**: 足球游戏的 Web 前端
- **proto**: 游戏和网关服务的 Protocol Buffer 定义

## 目录结构

- `game/go-server/cmd/main.go` - 游戏服务服务器实现
- `game/pkg/provider.go` - 游戏服务处理器实现
- `gate/go-server/cmd/main.go` - 带 HTTP 和 RPC 的网关服务服务器
- `gate/pkg/provider.go` - 网关服务处理器实现
- `gate/pkg/consumer.go` - 网关服务调用游戏服务的客户端
- `proto/` - Protocol Buffer 定义和生成的代码
- `website/` - Web 前端文件

## 运行方法

### 启动游戏服务

启动游戏服务服务器（监听端口 20000）：

```shell
go run ./game/go-server/cmd/main.go
```

### 启动网关服务

启动网关服务服务器（监听端口 20001 用于 RPC，8089 用于 HTTP）：

```shell
go run ./gate/go-server/cmd/main.go
```

### 访问 Web 前端

在浏览器中打开 `http://127.0.0.1:8089/` 访问游戏前端页面。前端将与网关服务的 HTTP API 通信。

## 服务通信

1. **前端 → 网关服务**: 对 `/login`、`/score`、`/rank` 的 HTTP 请求
2. **网关服务 → 游戏服务**: 使用 Triple 协议的 RPC 调用
3. **游戏服务 → 网关服务**: 用于网关操作的 RPC 调用

## 测试

您可以使用 curl 测试服务：

```shell
# 测试登录
curl "http://127.0.0.1:8089/login?name=player1"

# 测试得分
curl -X POST http://127.0.0.1:8089/score \
  -H "Content-Type: application/json" \
  -d '{"name":"player1","score":1}'

# 测试排名
curl "http://127.0.0.1:8089/rank?name=player1"
```

## 注意事项

- 在访问 Web 前端之前，请确保游戏服务器和网关服务器都在运行
- 游戏服务器必须在网关服务器之前启动，因为网关服务器需要连接到游戏服务

[version3]: https://protobuf.dev/programming-guides/proto3/
[Protocol Buffer Compiler Installation]: https://dubbo-next.staged.apache.org/zh-cn/overview/reference/protoc-installation/
