# Dubbo Golang 示例

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## 本工程包含的示例

* async: 异步调用
* attachment: 如何通过 attachment 把用户数据从调用方传递给服务方
* configcenter: 使用不同的配置中心，目前支持三种：zookeeper、apollo、和 nacos
* direct: 直连模式
* filter: 使用和扩展不同的 filter，目前包含了 custom_filter 和 tpslimit
* general: 通用例子，展示 zookeeper 注册中心的使用以及不同的配置项
* generic: 泛化调用
* helloworld: 入门例子
* multi_registry: 多注册
* registry: 展示与不同注册中心的对接，包含了 nacos、etcd、kubernetes 和 servicediscovery
* router: 展示了不同的路由，包含了 condition 和 tag
* seata: 展示了与 seata 的对接
* shop: 一个在线商店的小例子
* tracing: 链路追踪

## 如何运行

请参阅 [HOWTO](HOWTO_zh.md)

## 如何贡献

如果您希望增加新的用例，请继续阅读:

1. 为您的示例起合适的名字并创建子目录。如果您不太确定如何做，请参考现有示例摆放目录结构
2. 提交 PR 之前请确保在本地运行通过，提交 PR 之后请确保 GitHub 上的集成测试通过。请参考现有示例增加对应的测试
3. 请提供示例相关的 README.md 的中英文版本
