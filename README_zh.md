# Dubbo Golang 示例

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## 本工程包含的示例

* async: 异步调用
* attachment: 如何通过 attachment 把用户数据从调用方传递给服务方
* configcenter: 使用不同的配置中心，目前支持三种：zookeeper、apollo、和 nacos
* direct: 直连模式
* filter: 使用和扩展不同的 filter，目前包含了 custom_filter 和 tpslimit
* rpc: dubbo下展示dubbo协议通信调用
* generic: 泛化调用
* helloworld: 入门例子
* multi_registry: 多注册
* registry: 展示与不同注册中心的对接，包含了 nacos、etcd、kubernetes 和 servicediscovery
* router: 展示了 dubbo3 统一路由规则的使用，包括基于文件路由配置和k8s环境下的动态路由配置
* seata: 展示了与 seata 的对接
* shop: 一个在线商店的小例子
* tengine: 基于淘宝 Tengine 与 Dubbo-go 调用例子
* tls: 在 Dubbo-go 中使用 TLS
* tracing: 链路追踪
* game: 游戏服务例子

## 如何运行

请参阅 [HOWTO](HOWTO_zh.md)

## 如何贡献

如果您希望增加新的用例，请继续阅读:

1. 为您的示例起合适的名字并创建子目录。如果您不太确定如何做，请参考现有示例摆放目录结构
2. 提交 PR 之前请确保在本地运行通过，提交 PR 之后请确保 GitHub 上的集成测试通过。请参考现有示例增加对应的测试
3. 请提供示例相关的 README.md 的中英文版本
