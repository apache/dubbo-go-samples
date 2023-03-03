# Dubbo Go 示例仓库

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## 本工程包含的示例

* apisix: apache/apisix 与 Dubbo-go 整合示例
* async: dubbogo 通过 callback 方式进行异步 RPC 和call-one-way RPC
* config-api: 无需配置文件，使用 API 的方式启动服务、构造组件和使用
* configcenter: 使用不同的配置中心，目前支持：zookeeper、和 nacos
* context: 如何通过 go context  把用户数据 attachment 从调用方传递给服务方
* direct: 直连模式，无需注册中心，直连服务提供者
* error/triple: triple 示例，演示如何在 triple 协议中集成 hessian2/protobuf
* filter: dubbogo filter 示例，包含了 custom_filter 与 tpslimit
* game: dubbogo 游戏示例
* generic: 泛化调用
* helloworld: 入门例子
* integrate_test: dubbogo github action 集成测试用例
* logger: dubbogo logging
* mesh: 这个示例演示了如何将 dubbo-go 开发的应用部署在 Istio 体系下，以实现 Envoy 对 dubbo/dubbo-go 服务的自动代理
* metrics: 使用 prometheus 收集 dubbogo 的 metrics
* otel/trace: 使用 opentelemetry 进行 tracing
* proxyless: 这个示例演示了如何将 dubbo-go 开发的应用部署在 Istio 体系下，以实现 Envoy 对 dubbo/dubbo-go 服务的自动代理
* registry: 把 etcd/nacos/polaris/zookeeper 当做 dubbogo 注册中心示例
* rpc: 使用 Dubbogo 框架启动 rpc 服务，发起调用
  * rpc/dubbo: dubbo-go 3.0 RPC 通信示例，同时给出了 Java 和 Go 两种语言通信示例
  * rpc/grpc: 基于 gRPC 的 RPC 通信示例
  * rpc/jsonrpc: 基于json-rpc 的通信示例
  * rpc/triple: 基于 hessian2/msgpack/pb(protobuf-v3)/pb2(protobuf-v2)/自定义序列化协议 的序列化协议与 triple 通信协议相结合的 RPC 示例
* seata-go:  在 dubbogo 中如何基于 seata-go 实现分布式事务
* skywalking: 整合 skywalking 与 dubbogo 的示例
* tls: getty(tcp)/triple/gRPC 全链路 tls 安全通信示例
* tracing: 链路追踪例子，支持


## 如何运行

请参阅 [HOWTO](HOWTO_zh.md)

## 如何贡献

如果您希望增加新的用例，请继续阅读:

1. 为您的示例起合适的名字并创建子目录。如果您不太确定如何做，请参考现有示例摆放目录结构
2. 提交 PR 之前请确保在本地运行通过，提交 PR 之后请确保 GitHub 上的集成测试通过。请参考现有示例增加对应的测试
3. 请提供示例相关的 README.md 的中英文版本
* registry: 展示与不同注册中心的对接，包含了 nacos、etcd 和 zookeeper。
