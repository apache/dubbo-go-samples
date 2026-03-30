# 优雅停机示例

[English](README.md) | 中文

该示例用于验证 `dubbo-go` 中 Triple 协议的优雅停机流程。

它主要覆盖以下行为：

- 长连接消费者的主动通知
- 消费端在停机期间的被动关闭表现
- Provider 停机时对进行中请求的等待与排空
- `timeout`、`step-timeout`、`consumer-update-wait` 和 `offline-window` 等参数的影响

该示例**不包含注册中心**。因此你可以验证协议层的主动通知和请求排空行为，但不能直接观察“从注册中心摘除并传播”的完整链路。

## 前置条件

该示例使用仓库根目录下的 `go.mod`。

请在本地 `dubbo-go-samples` 仓库根目录执行命令：

```bash
cd /path/to/dubbo-go-samples
```

## 快速开始

在一个终端中启动服务端：

```bash
go run ./graceful_shutdown/go-server/cmd -timeout=60s -step-timeout=5s -delay=2s
```

在另一个终端中启动客户端：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=3 -interval=300ms -request-timeout=6s
```

然后在服务端终端按下 `Ctrl+C`。

预期现象：

- 已经在执行中的请求仍有机会完成
- 停机开始后，新请求会逐渐失败
- 服务端日志会按顺序输出优雅停机阶段

## 重要地址格式

进行直连调用时，`-addr` 必须带协议前缀。

例如：

- Triple：`tri://127.0.0.1:20000`

如果只传 `127.0.0.1:20000`，在某些场景下可能会被错误解析。

## 服务端参数

`graceful_shutdown/go-server/cmd/main.go` 支持以下参数：

- `-port=20000`
- `-timeout=60s`
- `-step-timeout=3s`
- `-consumer-update-wait=3s`
- `-offline-window=3s`
- `-delay=0s`

其中 `-delay` 会给每次请求增加固定处理延迟，用于观察停机时的在途请求排空效果。

## 客户端参数

`graceful_shutdown/go-client/cmd/main.go` 支持以下参数：

- `-addr=tri://127.0.0.1:20000`
- `-interval=200ms`
- `-concurrency=1`
- `-request-timeout=5s`
- `-short=true|false`
- `-name-prefix=hello`
- `-max-requests=0`
- `-min-successes=0`
- `-min-failures=0`

长连接验证时建议保持 `-short=false`。

其中 `-max-requests`、`-min-successes`、`-min-failures` 主要用于自动化测试。如果客户端在退出前没有达到这些最小阈值，会直接 `panic`，从而让集成测试失败。

## 推荐场景

### 1. 长连接下的 Triple 主动通知

终端 1：

```bash
go run ./graceful_shutdown/go-server/cmd -timeout=60s
```

终端 2：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=1 -interval=200ms
```

然后在服务端终端按下 `Ctrl+C`。

可观察点：

- 服务端日志会打印优雅停机各阶段
- 客户端在停机开始后不久会出现失败请求
- 长连接会收到主动通知，而不是仅在进程退出后断开

### 2. 在途请求排空

终端 1：

```bash
go run ./graceful_shutdown/go-server/cmd -delay=2s -step-timeout=5s
```

终端 2：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=3 -interval=300ms -request-timeout=6s
```

在请求尚未完成时于服务端终端按下 `Ctrl+C`。

可观察点：

- 已经开始执行的请求仍可能成功返回
- 新请求会在停机阶段开始失败
- 服务端会在等待预算耗尽或请求完成后退出

### 3. 同时观察主动通知与请求排空

缩短消费者更新等待时间，使停机更早开始拒绝新请求，同时保留已有请求排空窗口。

终端 1：

```bash
go run ./graceful_shutdown/go-server/cmd -delay=2s -timeout=15s -step-timeout=2s -consumer-update-wait=0s
```

终端 2：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -concurrency=2 -interval=200ms -request-timeout=4s
```

然后在服务端终端按下 `Ctrl+C`。

可观察点：

- 服务端日志会打印完整优雅停机序列
- 某些在途请求会在停机开始后继续完成
- 更新更晚到达的新请求会更早失败
- 客户端日志会体现 Triple 长连接的主动通知路径

### 4. 收紧整体超时预算

终端 1：

```bash
go run ./graceful_shutdown/go-server/cmd -timeout=10s -step-timeout=1s
```

终端 2：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000
```

该场景主要用于对比更紧整体优雅停机预算下的服务端日志表现。

### 5. 对比长连接与短连接

长连接：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000
```

短连接：

```bash
go run ./graceful_shutdown/go-client/cmd -addr=tri://127.0.0.1:20000 -short=true
```

针对主动通知路径，长连接更有代表性。

## 集成测试

该示例已接入根目录脚本驱动的集成测试：

```bash
./integrate_test.sh graceful_shutdown
```

脚本会启动 Triple 服务端，后台运行客户端，在观察到至少一次成功请求后向服务端发送中断信号，并要求客户端在退出前同时观察到：

- 至少一次成功请求
- 至少一次停机期间的失败请求

如果这些条件没有满足，客户端会直接 `panic`，从而使 CI 失败。

## 补充说明

- 该示例以 Triple 协议为主，用于聚焦当前优雅停机流程中的主动通知路径。
- 因为没有注册中心，这里只能覆盖协议层停机行为，不能完整覆盖注册中心摘除传播。
