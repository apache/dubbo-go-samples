# Dubbo Go & Polaris Router Example

[English](README.md) | 中文

## 使用服务路由功能

在 dubbogo 中快速体验北极星的服务路由能力

## 北极星服务端安装

[北极星服务端单机版本安装文档](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85/)

[北极星服务端集群版本安装文档](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E9%9B%86%E7%BE%A4%E7%89%88%E5%AE%89%E8%A3%85/)

## 如何使用

[北极星服务路由使用文档](https://polarismesh.cn/docs/%E5%8C%97%E6%9E%81%E6%98%9F%E6%98%AF%E4%BB%80%E4%B9%88/%E5%8A%9F%E8%83%BD%E7%89%B9%E6%80%A7/%E6%B5%81%E9%87%8F%E7%AE%A1%E7%90%86/#%E5%8A%A8%E6%80%81%E8%B7%AF%E7%94%B1)

### 如何配置服务路由参数

dubbogo 中的 PolarisMesh PriorityRouter 扩展点实现，能够根据用户配置的服务路由规则，自动的从当前 RPC 调用上下文以及请求信息中识别出需要参与服务路由的请求标签信息。

![](image/dubbogo-route-rule.png)

### 运行服务提供者

分别进入 server-prod、server-pre、server-dev 的 cmd 目录，执行以下命令：

```bash
# Dev 服务器 (端口 20000)
cd go-server/server-dev/cmd
go run .

# Pre 服务器 (端口 21000)
cd go-server/server-pre/cmd
go run .

# Prod 服务器 (端口 22000)
cd go-server/server-prod/cmd
go run .
```

当看到以下日志时即表示 server 端启动成功：

```
INFO ... dubbo server started
```

### 运行服务调用者

进入 go-client 的 cmd 目录，执行以下命令：

```bash
cd go-client/cmd

# 无 uid (路由到 prod)
export uid=
go run .

# uid=user-1 (路由到 pre)
export uid=user-1
go run .

# uid=user-2 (路由到 dev)
export uid=user-2
go run .
```

当看到以下日志时即表示 go-client 成功发现 go-server 并发起了 RPC 调用：

```
INFO ... uid=, response: user:"[Prod] Alex Stocks"
```

### 预期输出

当路由正确工作时，您将看到：

```
# export uid=
uid=, response: user:"[Prod] Alex Stocks"
uid=, response: user:"[Prod] Alex Stocks"

# export uid=user-1
uid=user-1, response: user:"[Pre] Alex Stocks"
uid=user-1, response: user:"[Pre] Alex Stocks"

# export uid=user-2
uid=user-2, response: user:"[Dev] Alex Stocks"
uid=user-2, response: user:"[Dev] Alex Stocks"
```
