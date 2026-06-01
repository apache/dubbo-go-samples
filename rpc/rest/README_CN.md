# Dubbo-Go REST 示例

这个示例用于验证 Dubbo-Go REST 协议的三种地址获取方式：

- 直连 URL
- 接口级注册
- 应用级服务发现

示例通过本地 REST 配置显式声明 HTTP 方法、路径参数、查询参数、请求头和请求体映射。

## 运行

默认模式是 Nacos 应用级服务发现：

```bash
go run ./rpc/rest/go-server/cmd
```

在另一个终端运行消费者：

```bash
go run ./rpc/rest/go-client/cmd
```

期望输出：

```text
REST response: userID=101 name=dubbo-go traceID=trace-rest-basic message=body-from-dubbo-rest-client greeting="hello dubbo-go, userID=101, traceID=trace-rest-basic, message=body-from-dubbo-rest-client"
```

Provider 同时也是一个普通 HTTP 服务，可以直接用 `curl` 调用：

```bash
curl -s \
  -X POST 'http://127.0.0.1:20080/api/v1/users/202/greeting?name=curl' \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'X-Trace-ID: trace-curl' \
  -d '{"message":"body-from-curl"}'
```

## REST 映射

Provider URL 只提供网络目标，例如：

```text
rest://127.0.0.1:20080/org.apache.dubbo.samples.rest.GreetingService
```

真正的 REST 调用形态来自 `api/rest_config.go`：

- HTTP 方法：`POST`
- 路径：`/api/v1/users/{userID}/greeting`
- 路径参数：参数 `0` -> `userID`
- 查询参数：参数 `1` -> `name`
- 请求头：参数 `2` -> `X-Trace-ID`
- 请求体：参数 `3`

也就是说，这个示例刻意区分了“服务地址发现”和“REST HTTP 映射”。

## 服务发现模式

Provider 和 Consumer 都支持以下参数：

- `-registry=direct|zookeeper|nacos`
- `-registry-type=interface|service|all`

默认值是：

```bash
-registry=nacos -registry-type=service
```

`interface` 表示注册中心直接保存可调用的 `rest://...` provider URL。  
`service` 表示注册中心保存应用实例；consumer 先通过 service mapping 找到应用，再通过 metadata service 获取服务元数据，最后根据应用实例和 `ServiceInfo` 重建 `rest://...` provider URL。  
`all` 表示同时注册接口级和应用级两种数据。

## 直连模式

```bash
go run ./rpc/rest/go-server/cmd -registry=direct
go run ./rpc/rest/go-client/cmd -registry=direct
```

## ZooKeeper 接口级注册

```bash
go run ./rpc/rest/go-server/cmd -registry=zookeeper -registry-type=interface
go run ./rpc/rest/go-client/cmd -registry=zookeeper -registry-type=interface
```

所有模式下，consumer 都应该输出相同的 `REST response: ...`。
