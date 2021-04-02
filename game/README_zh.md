# 游戏服务示例

### 背景介绍

- 示例包含 **gate** (网关服务) 和 **game** (逻辑服务) 两个服务
- 两个服务会互相 RPC 通讯 (都同时注册 **provider** 和 **consumer**)
- **gate** 额外启动了 http 服务 (端口 **8000**), 用于手工触发 **gate** RPC调用 **game**

> 每次 **gate** RPC调用(**Message**) **game** 后, **game** 会同步RPC调用(Send) **gate** 推送相同消息
> 此逻辑默认已注释, 位置在 `go-server-game/pkg/provider.go` (25 ~ 30行)


请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。
