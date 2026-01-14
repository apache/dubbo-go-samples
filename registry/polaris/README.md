# Polaris as registry

[English](README.md) | [中文](README_zh.md)

This example shows dubbo-go's service discovery feature with Polaris as registry.

## How to run

### Start Polaris server
Follow this instruction to [install and start Polaris server](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85/).

### Run server
```shell
$ go run ./go-server/cmd/server.go
```

Open Polaris console, check service successfully registered into Polaris.

### Run client
```shell
$ go run ./go-client/cmd/client.go
```
