# Dubbo Golang Examples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

[中文 🇨🇳](./README_CN.md)

## What It Contains

* apisix: apache/apisix and Dubbo-go example
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
* otel/trace: show how to use opentelemetry as dubbogo tracing tool
* proxyless: give an proxyless service mesh example which shows how to deploy the dubbo-go services on kubernetes platform
* registry: show how to use etcd/nacos/polaris/zookeeper as dubbogo registry
* rpc: dubbo directory display dubbo protocol communication
* seata: A seata example
* skywalking: show how to integrate skywalking into dubbogo
* tracing: tracing example

## How To Run

Pls refer [How To Run](HOWTO.md) for the instructions.

## How to contribute

If you want to add more samples, pls. read on:
1. Create new sub directory and give it an appropriate name for your new sample. Pls follow the layout of the existing sample if you are not sure how to organize your code.
2. Make sure your sample work as expected before submit PR, and make sure GitHub CI passes after PR is submitted. Pls refer to the existing sample on how to test the sample.
3. Pls provide README.md to explain your samples.
