# Dubbo Go 示例仓库

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## 本工程包含的示例

* config-api: 无需配置文件，使用 API 的方式启动服务、构造组件和使用。
* configcenter: 使用不同的配置中心，目前支持：zookeeper、apollo、和 nacos
* context: 如何通过 go context  把用户数据 attachment 从调用方传递给服务方
* direct: 直连模式，无需注册中心，直连服务提供者
* rpc: 使用 Dubbogo 框架启动 rpc 服务，发起调用
* generic: 泛化调用
* helloworld: 入门例子
* registry: 展示与不同注册中心的对接，包含了 nacos、etcd 和 zookeeper。
* tracing: 链路追踪例子，支持

## 如何运行

请参阅 [HOWTO](HOWTO_zh.md)

## 如何贡献

如果您希望增加新的用例，请继续阅读:

1. 为您的示例起合适的名字并创建子目录。如果您不太确定如何做，请参考现有示例摆放目录结构
2. 提交 PR 之前请确保在本地运行通过，提交 PR 之后请确保 GitHub 上的集成测试通过。请参考现有示例增加对应的测试
3. 请提供示例相关的 README.md 的中英文版本
