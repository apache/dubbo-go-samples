# dubbo-go 标签路由示例

## 1 标签路由介绍

标签路由是 dubbo 流量管控的重要部分。标签路由通过将某一个服务的实例划分到不同的分组，约束具有特定标签的流量只能在指定分组中流转，不同分组为不同的流量场景服务，从而达到实现流量隔离的目的，可以作为蓝绿发布、灰度发布等场景能力的基础。

## 2 使用方法
### 2.1 provider

在 provider 端需要对服务的提供者打标签，目前有两种方式：静态打标、动态打标。

静态打标是在服务启动之前在配置文件中对服务提供者打标签。每次更换标签都需要重启服务提供者。

如下图所示，对服务打了 tag1 的标签。**注意：打标签的对象是应用级别。**
```yaml
dubbo:
  application:
    tag: tag1
  registries:
    demoZK:
      protocol: zookeeper
      address: 127.0.0.1:2181
  protocols:
    triple:
      name: tri
      port: 20000
  provider:
    services:
      GreeterProvider:
        interface: com.apache.dubbo.sample.basic.IGreeter # must be compatible with grpc or dubbo-java
```
动态打标相对于静态打标更加灵活，可以在服务提供的过程中对服务的标签进行更换。但是需要用户通过 dubbo-admin 操作。

进入 dubbo-admin 服务治理下的标签路由模块。点击创建会弹出创建标签路由窗口。在规则内容里面可以配置相关规则。
```yaml
force: false
enabled: true
runtime: false
tags:
  - name: tag1
    addresses: [192.168.0.1:20881]
  - name: tag2
    addresses: [192.168.0.2:20882]
```
配置完成后点击保存就可以动态修改服务提供者的标签。 

### 2.2 consumer

在 consumer 端可以选择使用哪个标签的服务提供商，主要在代码中进行定义。如下是使用 tag1 提供的服务。
```go
// set tag
ctx := context.Background()
atm := map[string]string{
    "dubbo.tag":       "tag1",
    "dubbo.force.tag": "true",
}
ctx = context.WithValue(ctx, constant.AttachmentKey, atm)
reply, err := grpcGreeterImpl.SayHello(ctx, req)
if err != nil {
    logger.Error(err)
}
logger.Infof("client response result: %v\n", reply)
```

## 3 部署

用户可以部署该 demo 尝试使用标签路由的功能，dubbo-go 需要升级到最新版本。
1. 启动注册中心 zookeeper 和 dubbo-admin；
2. 设置 provider 的配置文件环境变量并启动 provider；
3. 设置 consumer 的配置文件环境变量并启动 consumer，可以看到消费 provider 的服务；
4. 修改 provider 的 tag 后重启尝试；
5. 通过配置中心修改 provider 的 tag 后尝试。