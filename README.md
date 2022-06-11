# Dubbo Golang Examples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## What It Contains

* async: An async example.
* attachment: An attachment example, to show how to use attachment to pass user data from the client to the server.
* configcenter: uses of different config centers, including zookeeper and nacos at present.
* direct: A direct invocation example.
* filter: Some examples of different filter, including custom_filter and tpslimit
* rpc: dubbo directory display dubbo protocol communication
* generic: A generic invocation example
* helloworld: A 101 example
* router: router examples, including condition and tag
* seata: A seata example
* shop: Shop sample
* tengine: Taobao Tengine and Dubbo-go example
* tls: Use TLS in Dubbo-go application
* tracing: tracing example
* game: game service example

## How To Run

Pls refer [How To Run](HOWTO.md) for the instructions.

## How to contribute

If you want to add more samples, pls. read on:
1. Create new sub directory and give it an appropriate name for your new sample. Pls follow the layout of the existing sample if you are not sure how to organize your code.
2. Make sure your sample work as expected before submit PR, and make sure GitHub CI passes after PR is submitted. Pls refer to the existing sample on how to test the sample.
3. Pls provide README.md to explain your samples.
