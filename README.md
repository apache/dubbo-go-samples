# Dubbo Golang Examples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

[中文 🇨🇳](./README_CN.md)

## What It Contains

* apisix: apache/apisix and Dubbo-go example
* async: the dubbogo callback[asynchronous] and call-one-way rpc example
* config-api: show how to use dubbogo by APIs without configuration
* configcenter: uses of different config centers, including zookeeper and nacos at present.
* context: how to transfer request context between multiple producers/consumers
* direct: A direct invocation example.
* error/triple: triple sample with hessian2/protobuf
* filter: Some examples of different filter, including custom_filter and tpslimit
* game: game service example
* generic: A generic invocation example
* helloworld: A 101 example
* integrate_test: dubbogo github action integrate test use cases
* logger: dubbogo logging
* mesh: give an proxy service mesh example which shows how to deploy the dubbo-go services with envoy on kubernetes platform
* metrics: show how to collect dubbogo prometheus metrics
* multirpc: show how to use three protocols(triple, dubbo, jsonrpc) in one server and invoke them on the client-side
* otel/trace: show how to use opentelemetry as dubbogo tracing tool
* proxyless: give an proxyless service mesh example which shows how to deploy the dubbo-go services on kubernetes platform
* registry: show how to use etcd/nacos/polaris/zookeeper as dubbogo registry
* rpc: dubbo directory display dubbo protocol communication
  * rpc/dubbo: dubbo-go 3.0 rpc example between Java and Go
  * rpc/grpc: dubbo-go rpc example based on gRPC
  * rpc/jsonrpc: dubbo-go rpc example based on json-rpc
  * rpc/triple: dubbo-go rpc example based on hessian2/msgpack/pb(protobuf-v3)/pb2(protobuf-v2)/self-defined-serialization
* seata-go: A seata-go example
* skywalking: show how to integrate skywalking into dubbogo
* tls: use TLS encryption in getty(tcp)/triple/gRPC communication mode
* tracing: tracing example
* error: error handling
* llm: dubbo-go integration with large language models (LLM), using Ollama for inference

## How To Run

Pls refer [How To Run](HOWTO.md) for the instructions.

## How to contribute

If you want to add more samples, pls. read on:
1. Create new sub directory and give it an appropriate name for your new sample. Pls follow the layout of the existing sample if you are not sure how to organize your code.
2. Make sure your sample work as expected before submit PR, and make sure GitHub CI passes after PR is submitted. Pls refer to the existing sample on how to test the sample.
3. Pls provide README.md to explain your samples.
