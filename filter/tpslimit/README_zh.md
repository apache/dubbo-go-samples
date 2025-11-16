# TPSLimit Filter 示例

[English](README.md) | [中文](README_zh.md)

### 背景

Dubbo-go 内置了限流 filter "tpslimit"。可以通过在服务端的配置来激活，另外，用户还可以自定义限流策略和拒绝访问后的处理逻辑。

### 示例

##### 1. 代码

A) 自定义限流策略：

通过实现 filter.TpsLimitStrategy 来自定义限流策略。在本例中，采取的策略是随机限流。例子链接为：[limit_strategy.go](go-server/pkg/limit_strategy.go)

B) 自定义拒绝访问处理：

通过实现 filter.RejectedExecutionHandler。在本例中，当限流条件满足的情况下，拒绝访问的自定义处理逻辑将会返回 "The request is rejected and doesn't have any default value. " 的错误给客户端。例子链接为：[reject_handler.go](go-server/pkg/reject_handler.go)

##### 2. 配置

在服务端代码中使用 v3 API 启用 tpslimit filter：

```go
if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{},
	server.WithTpsLimiter("method-service"),
); err != nil {
	panic(err)
}
```

自定义的 TPS 限流策略和拒绝执行处理器通过各自包的 `init()` 函数中的 `extension.SetTpsLimitStrategy()` 和 `extension.SetRejectedExecutionHandler()` 进行注册。

##### 3. 运行

### 前置条件

1. 启动 Zookeeper（默认地址：`127.0.0.1:2181`）

### 运行服务端

```shell
go run ./go-server/cmd/main.go
```

### 运行客户端

```shell
go run ./go-client/cmd/main.go
```

## 预期输出

**服务端输出：**

```bash
Random IsAllowable!
ERROR   The invocation was rejected due to over the tps limitation, ...
```

**客户端输出：**

```bash
start to test tpslimit
error: The request is rejected and doesn't have any default value.
response: hello world
...
successCount=<数量>, failCount=<数量>
```

客户端将发送 60 个请求，每个请求间隔 200ms。部分请求会被 TPS 限流器拒绝，最终输出中会显示成功和失败的请求数量。

