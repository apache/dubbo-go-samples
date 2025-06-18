# Dubbo Go 示例仓库

## 本工程包含的示例

### compatibility（旧版 dubbo-go 示例，部分已迁移，推荐优先参考此目录）
* compatibility/apisix：apache/apisix 与 Dubbo-go 整合示例
* compatibility/async：通过 callback 方式进行异步 RPC 及单向调用示例
* compatibility/config-api：无需配置文件，使用 API 启动 Dubbo-go 服务
* compatibility/configcenter：多种配置中心（如 zookeeper、nacos）用法示例
* compatibility/context：多生产者/消费者间传递请求 context 示例
* compatibility/direct：直连调用示例，无需注册中心
* compatibility/error：错误处理与 triple 协议示例
* compatibility/filter：内置及自定义 filter 示例（如 custom_filter、tpslimit）
* compatibility/game：游戏服务示例
* compatibility/generic：泛化调用示例
* compatibility/logger：Dubbo-go 日志功能示例
* compatibility/mesh：基于 Envoy 的服务网格部署示例
* compatibility/metrics：Prometheus 指标采集示例
* compatibility/otel：OpenTelemetry 链路追踪示例
* compatibility/proxyless：Kubernetes 下 proxyless 服务网格示例
* compatibility/registry：多种注册中心（etcd/nacos/polaris/zookeeper）用法
* compatibility/rpc：Dubbo 协议通信示例
  * compatibility/rpc/dubbo：Dubbo-go 3.0 Java/Go 跨语言 RPC 示例
  * compatibility/rpc/grpc：基于 gRPC 的 Dubbo-go RPC 示例
  * compatibility/rpc/jsonrpc：基于 JSON-RPC 的 Dubbo-go RPC 示例
  * compatibility/rpc/triple：支持多种序列化（hessian2/msgpack/pb/pb2/自定义）的 triple 协议示例
* compatibility/seata-go：seata-go 分布式事务示例
* compatibility/skywalking：集成 SkyWalking 的 Dubbo-go 示例
* compatibility/tls：getty(tcp)/triple/gRPC 全链路 TLS 加密通信示例
* compatibility/tracing：链路追踪示例

### 传统示例（部分已迁移至 compatibility 目录，建议优先参考 compatibility）
* book-flight-ai-agent：AI agent 机票预订示例
* config_center：多种配置中心（如 nacos、zookeeper）管理配置示例
* config_yaml：通过 YAML 文件配置 Dubbo-go 应用
* context：通过 go context 传递用户数据（attachments）的示例
* error：Dubbo-go 错误处理示例
* filter：内置及自定义 filter 用法示例
* healthcheck：服务健康检查示例
* helloworld：Dubbo-go 入门 Hello World 示例
* integrate_test：Dubbo-go 示例集成测试用例
* java_interop：Java 与 Go Dubbo 实现互操作示例
* llm：集成大语言模型（LLM）与 Dubbo-go 示例
* logger：Dubbo-go 日志功能示例
* metrics：Dubbo-go 服务指标采集与暴露示例
* online_boutique：Dubbo-go 微服务电商示例
* otel/tracing：基于 OpenTelemetry 的分布式链路追踪示例
* registry：多种服务注册中心（如 nacos、zookeeper）用法示例
* retry：Dubbo-go RPC 调用重试机制示例
* rpc：多种 RPC 协议示例
  * rpc/dubbo：Dubbo 协议示例，含 Java/Go 跨语言通信
  * rpc/grpc：gRPC 协议示例
  * rpc/jsonrpc：JSON-RPC 协议示例
  * rpc/triple：支持多种序列化的 triple 协议示例
* streaming：流式 RPC 通信示例
* task：任务调度与执行示例
* timeout：Dubbo-go 超时处理示例
* transaction/seata-go：基于 seata-go 的分布式事务示例

## 如何运行

请参考 [How To Run](HOWTO_zh.md) 获取详细运行说明。

## 如何贡献

如需新增示例，请遵循以下流程：
1. 新建子目录并命名，建议参考现有示例目录结构。
2. 确保示例可正常运行，提交 PR 后通过 GitHub CI 检查。可参考现有示例的测试方式。
3. 请为你的示例提供 README.md 说明文档。
* registry: 展示与不同注册中心的对接，包含了 nacos、etcd 和 zookeeper。
