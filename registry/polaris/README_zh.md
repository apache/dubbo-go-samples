# 使用 Polaris 作为注册中心

[English](README.md) | [中文](README_zh.md)

本示例展示 dubbo-go 使用 Polaris 作为注册中心的服务发现功能。

## 如何运行

### 启动 Polaris 服务器
按照[安装文档](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85/)安装并启动 Polaris 服务器。

### 运行服务端
```shell
$ go run ./go-server/cmd/server.go
```

打开 Polaris 控制台，检查服务是否成功注册到 Polaris。

### 运行客户端
```shell
$ go run ./go-client/cmd/client.go
```

