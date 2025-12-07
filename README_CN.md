# Dubbo-Go Samples 

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

中文 | [English](README.md)

一组可运行的 Dubbo-go 示例，涵盖配置、注册中心、可观测性、互操作性、服务网格等场景。

## 内容概览

### 示例

* `async`：回调（异步）与单向 RPC 调用示例。
* `book-flight-ai-agent`：使用 AI Agent 实现机票预订的示例。
* `config_center`：演示如何使用不同配置中心（如 Nacos、Zookeeper）进行配置管理。
* `config_yaml`：展示如何使用 YAML 文件配置 Dubbo-go 应用。
* `context`：演示如何通过 Go 的 `context` 在客户端与服务端之间传递用户数据（attachments）。
* `error`：Dubbo-go 的错误处理示例。
* `filter`：演示 Dubbo-go 中内置和自定义 Filter 的使用。
  * `polaris/limit`：使用 Polaris 实现 TPS 限流。
* `healthcheck`：服务健康检查示例。
* `helloworld`：Dubbo-go 最基础的 “Hello World” 示例，同时包含 Go 与 Java 的互操作示例。
* `http3`：HTTP/3（QUIC）协议支持示例，演示如何通过 Triple 协议使用 HTTP/3 实现 Go 与 Java 服务之间的高性能通信，并支持 TLS 加密。
* `direct`：不依赖注册中心的 Triple 点对点调用示例。
* `game`：游戏服务示例。
* `generic`：泛化调用示例，支持 Dubbo-Go 与 Dubbo Java 服务互操作，适用于无接口信息场景。
* `integrate_test`：Dubbo-go 示例的集成测试用例。
* `java_interop`：展示 Java 与 Go Dubbo 实现之间的互操作能力。
* `llm`：将大模型（LLM）集成到 Dubbo-go 中的示例。
* `logger`：Dubbo-go 应用的日志使用示例。
* `metrics`：展示如何采集并暴露 Dubbo-go 服务的指标，支持 Prometheus Push 和 Pull 两种模式；同时包含用于清理 Push 模式僵尸指标的 `pgw-cleaner` 工具。
* `online_boutique`：基于 Dubbo-go 构建的微服务 “在线商城” 演示项目。
* `otel/tracing`：使用 OpenTelemetry 的分布式链路追踪示例。
  * `stdout`：使用 stdout exporter 输出追踪数据。
  * `otlp_http_exporter`：使用 `otlpHttpExporter` 导出追踪数据，覆盖 `dubbo` / `triple` / `jsonrpc` 协议。
* `registry`：使用不同服务注册中心（如 Nacos、Zookeeper）的示例。
* `retry`：Dubbo-go RPC 调用重试机制示例。
* `rpc`：Dubbo-go 支持的多种 RPC 协议示例。
  * `rpc/dubbo`：Dubbo 协议示例，包含 Java–Go 互操作。
  * `rpc/grpc`：基于 gRPC 协议的示例。
  * `rpc/jsonrpc`：基于 JSON-RPC 协议的示例。
  * `rpc/triple`：Triple 协议示例，涵盖多种序列化方式。
* `streaming`：流式 RPC 调用示例。
* `task`：任务调度与执行示例。
* `timeout`：Dubbo-go 超时处理示例。
* `tls`：演示如何在 Dubbo-go 中使用 TLS（基于 X.509 证书），实现客户端与服务端之间的加密通信和/或双向认证。
* `transaction/seata-go`：基于 `seata-go` 的分布式事务示例。

### compatibility（旧版 Dubbo-go 示例）

* `compatibility/apisix`：Dubbo-go 集成 Apache APISIX 的示例。
* `compatibility/config-api`：演示如何在不使用配置文件的情况下，通过 API 使用 Dubbo-go。
* `compatibility/configcenter`：不同配置中心的使用示例，包括 Zookeeper、Nacos 和 Apollo。
* `compatibility/generic`：泛化调用示例。
* `compatibility/mesh`：基于代理的服务网格示例，展示如何在 Kubernetes 上结合 Envoy 部署 Dubbo-go 服务。
* `compatibility/proxyless`：无 Sidecar 的服务网格示例，展示在 Kubernetes 上的部署方式。
* `compatibility/polaris`:  Dubbo-go 与 polaris 集成示例.
    * `compatibility/polaris/limit`: 在 dubbogo 中快速体验北极星的服务限流能力
    * `compatibility/polaris/registry`: 在 dubbogo 中快速体验北极星的服务注册以及服务发现能力
    * `compatibility/polaris/router`: 在 dubbogo 中快速体验北极星的服务路由能力
* `compatibility/registry`：演示将 `etcd` / `Nacos` / `Zookeeper` 作为 Dubbo-go 注册中心的使用方法。
* `compatibility/rpc`：Dubbo 协议通信示例。
  * `compatibility/rpc/dubbo`：Dubbo-go 3.0 下 Java 与 Go 的 RPC 示例。
  * `compatibility/rpc/grpc`：基于 gRPC 的 Dubbo-go RPC 示例。
  * `compatibility/rpc/jsonrpc`：基于 JSON-RPC 的 Dubbo-go RPC 示例。
  * `compatibility/rpc/triple`：基于 `hessian2` / `msgpack` / `pb`（protobuf v3）/ `pb2`（protobuf v2）/ 自定义序列化 的 Triple 协议示例。
* `compatibility/skywalking`：如何在 Dubbo-go 中集成 SkyWalking。
* `compatibility/tls`：在 `getty`（TCP）/`triple`/`gRPC` 通信模式下使用 TLS 加密的示例。
* `compatibility/tracing`：链路追踪示例。

### 工具

* `pgw-cleaner`：用于在 Prometheus Push 模式下清理僵尸指标的运维工具。

## 运行示例

请参考 [HOWTO.md](HOWTO.md) 获取运行各个示例的详细说明。

## 如何参与贡献

如果你希望添加更多示例，请按以下步骤进行：

1. 新建一个子目录，并为你的示例取一个合适的名称。如果不确定如何组织代码，可以参考现有示例的目录结构。
2. 在提交 PR 之前，请确保示例能正常运行；提交 PR 后，也请确保 GitHub CI 能通过。可以参考已有示例了解如何编写和执行测试。
3. 在你的示例目录下提供一个 `README.md`，说明该示例的功能及运行方式。
