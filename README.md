# Dubbo Golang Examples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

[中文 🇨🇳](./README_CN.md)

## What It Contains

### compatibility (old dubbo-go samples)
* compatibility/apisix: Example integrating apache/apisix with Dubbo-go
* compatibility/async: Callback (asynchronous) and one-way RPC example
* compatibility/config-api: How to use Dubbo-go by APIs without configuration files
* compatibility/configcenter: Usage of different config centers, including zookeeper and nacos
* compatibility/context: How to transfer request context between multiple producers/consumers
* compatibility/direct: Direct invocation example
* compatibility/error: Error handling and triple protocol examples
* compatibility/filter: Examples of different filters, including custom_filter and tpslimit
* compatibility/game: Game service example
* compatibility/generic: Generic invocation example
* compatibility/logger: Dubbo-go logging examples
* compatibility/mesh: Proxy service mesh example showing how to deploy Dubbo-go services with Envoy on Kubernetes
* compatibility/metrics: How to collect Dubbo-go Prometheus metrics
* compatibility/otel: How to use OpenTelemetry as Dubbo-go tracing tool
* compatibility/proxyless: Proxyless service mesh example for deploying Dubbo-go services on Kubernetes
* compatibility/registry: How to use etcd/nacos/polaris/zookeeper as Dubbo-go registry
* compatibility/rpc: Dubbo protocol communication examples
  * compatibility/rpc/dubbo: Dubbo-go 3.0 RPC example between Java and Go
  * compatibility/rpc/grpc: Dubbo-go RPC example based on gRPC
  * compatibility/rpc/jsonrpc: Dubbo-go RPC example based on JSON-RPC
  * compatibility/rpc/triple: Dubbo-go RPC example based on hessian2/msgpack/pb(protobuf-v3)/pb2(protobuf-v2)/self-defined-serialization
* compatibility/seata-go: Seata-go distributed transaction example
* compatibility/skywalking: How to integrate SkyWalking into Dubbo-go
* compatibility/tls: Use TLS encryption in getty (tcp)/triple/gRPC communication mode
* compatibility/tracing: Tracing example

### Legacy samples 
* book-flight-ai-agent: Example for booking flights using an AI agent
* config_center: Demonstrates how to use different config centers (e.g., nacos, zookeeper) for configuration management
* config_yaml: Shows how to configure Dubbo-go applications using YAML files
* context: Example of passing user data (attachments) via Go context between client and server
* error: Error handling examples in Dubbo-go
* filter: Demonstrates the use of built-in and custom filters in Dubbo-go
* healthcheck: Service health check example
* helloworld: Basic hello world example for Dubbo-go
* integrate_test: Integration test cases for Dubbo-go samples
* java_interop: Demonstrates interoperability between Java and Go Dubbo implementations
* llm: Example of integrating large language models (LLM) with Dubbo-go
* logger: Logging examples for Dubbo-go applications
* metrics: How to collect and expose metrics from Dubbo-go services
* online_boutique: Microservices online boutique demo using Dubbo-go
* otel/tracing: Distributed tracing example using OpenTelemetry
* registry: Examples of using different service registries (e.g., nacos, zookeeper)
* retry: Demonstrates retry mechanisms in Dubbo-go RPC calls
* rpc: Various RPC protocol examples with Dubbo-go
  * rpc/dubbo: Dubbo protocol example, including Java and Go interop
  * rpc/grpc: gRPC protocol example
  * rpc/jsonrpc: JSON-RPC protocol example
  * rpc/triple: Triple protocol example with multiple serialization formats
* streaming: Streaming RPC example
* task: Task scheduling and execution example
* timeout: Demonstrates timeout handling in Dubbo-go
* transaction/seata-go: Distributed transaction example using seata-go

## How To Run

Please refer to [How To Run](HOWTO.md) for instructions.

## How to contribute

If you want to add more samples, please read on:
1. Create a new subdirectory and give it an appropriate name for your new sample. Please follow the layout of the existing samples if you are not sure how to organize your code.
2. Make sure your sample works as expected before submitting a PR, and ensure GitHub CI passes after the PR is submitted. Please refer to the existing samples on how to test your sample.
3. Please provide a README.md to explain your sample.
