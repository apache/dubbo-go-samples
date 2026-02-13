# Dubbo-Go Samples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

English | [中文](README_CN.md)

A collection of runnable Dubbo-go examples covering configuration, registries, observability, interoperability, service mesh, and more.

## What It Contains

### Samples

* `async`: Callback (asynchronous) and one-way RPC examples.
* `apisix`: Example integrating Apache APISIX with Dubbo-go.
* `book-flight-ai-agent`: Example of booking flights using an AI agent.
* `config_center`: Demonstrates how to use different config centers (e.g., Nacos, Zookeeper) for configuration management.
* `config_yaml`: Shows how to configure Dubbo-go applications using YAML files.
* `context`: Demonstrates passing user data (attachments) via Go `context` between client and server.
* `error`: Error-handling examples in Dubbo-go.
* `filter`: Demonstrates the use of built-in and custom filters in Dubbo-go.
  * `polaris/limit`: Uses Polaris as a TPS limiter.
* `healthcheck`: Service health check example.
* `helloworld`: Basic “Hello World” example for Dubbo-go, also includes Go–Java interoperability.
* `http3`: HTTP/3 (QUIC) protocol support example demonstrating how to use Triple protocol with HTTP/3 for high-performance communication between Go and Java services with TLS encryption.
* `direct`: Triple point-to-point invocation sample without a registry, also includes Go–Java interoperability.
* `game`: Game service example.
* `generic`: Generic invocation examples supporting interoperability between Dubbo-Go and Dubbo Java services, suitable for scenarios without interface information.
* `integrate_test`: Integration test cases for Dubbo-go samples.
* `java_interop`: Demonstrates interoperability between Java and Go Dubbo implementations.
* `llm`: Example of integrating Large Language Models (LLMs) with Dubbo-go.
* `logger`: Logging examples for Dubbo-go applications.
* `metrics`: Observability-related samples.
  * `metrics/prometheus_grafana`: Shows how to collect and expose metrics from Dubbo-go services, supporting both Prometheus Push and Pull modes. Also includes the `pgw-cleaner` tool for cleaning zombie metrics in Push mode.
  * `metrics/probe`: Demonstrates Dubbo-go Kubernetes probe endpoints (`/live`, `/ready`, `/startup`) and deployment usage.
* `mesh`: Proxy-based service mesh example showing how to deploy Dubbo-go services with Envoy on Kubernetes.
* `online_boutique`: Microservices “online boutique” demo built with Dubbo-go.
* `otel/tracing`: Distributed tracing examples using OpenTelemetry.
  * `stdout`: Uses stdout exporter to export tracing data.
  * `otlp_http_exporter`: Uses `otlpHttpExporter` to export tracing data, covering `dubbo`/`triple`/`jsonrpc` protocols.
* `registry`: Examples of using different service registries (e.g., Nacos, Zookeeper).
* `retry`: Demonstrates retry mechanisms in Dubbo-go RPC calls.
* `rpc`: Various RPC protocol examples with Dubbo-go.
  * `rpc/dubbo`: Dubbo protocol example, including Java–Go interop.
  * `rpc/grpc`: gRPC protocol example.
  * `rpc/jsonrpc`: JSON-RPC protocol example.
  * `rpc/triple`: Triple protocol example with multiple serialization formats.
* `streaming`: Streaming RPC example, also includes Go–Java interoperability when both use streaming.
* `task`: Task scheduling and execution example.
* `timeout`: Demonstrates timeout handling in Dubbo-go.
* `tls`: Demonstrates how to use TLS (based on X.509 certificates) in Dubbo-go to enable encrypted communication and/or mutual authentication between client and server, also includes Go–Java interoperability.
* `transaction/seata-go`: Distributed transaction example using `seata-go`.

### compatibility (legacy Dubbo-go samples)

* `compatibility/generic`: Generic invocation example.
* `compatibility/polaris`:  Dubbo-go integrate with polaris samples.
    * `compatibility/polaris/router`: Quickly experience Polaris' service routing capabilities in dubbogo

### Tools

* `pgw-cleaner`: Operations and maintenance tool for cleaning up zombie metrics in Prometheus Push mode.

## How To Run

Please refer to [HOWTO.md](HOWTO.md) for detailed instructions on running the samples.

## How to Contribute

If you want to add more samples, please follow these steps:

1. Create a new subdirectory and give it an appropriate name for your sample. If you are unsure how to organize your code, follow the layout of the existing samples.
2. Make sure your sample works as expected before submitting a PR, and ensure GitHub CI passes after the PR is submitted. You can refer to the existing samples to learn how to test your sample.
3. Provide a `README.md` in your sample directory to explain what your sample does and how to run it.
