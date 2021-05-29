# Dubbo Golang Examples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## What It Contains

* async: An async example.
* attachment: An attachment example, to show how to use attachment to pass user data from the client to the server.  
* chain: Show an example of calling chain
* config-api: Use API to config Dubbo-go
* configcenter: Uses of different config centers, including zookeeper, apollo and nacos at present.
* context: Use context
* direct: A direct invocation example
* docker: Use Docker to package and run Dubbo-go application
* filter: Some examples of different filter, including custom_filter and tpslimit
* game: game service example
* general: A general example
* generic: A generic invocation example
* group: Service group
* helloworld: A 101 example
* metric: Enable metrics in Dubbo-go application
* multi-registry: Register services to different registry centers
* multi-zone: Call Dubbo services from different zones
* multi_registry: A multi-registry example
* registry: uses of different registres, including kubernetes, nacos, etcd and service-discovery
* router: router examples, including condition and tag
* seata: A seata example
* shopping-center: A complete shopping sample
* tengine: Tengine Calls Dubbo-go example
* tls: Use TLS in Dubbo-go application
* tracing: tracing example

## How To Run

Pls. refer [How To Run](HOWTO.md) for the instructions.

## How to contribute

If you want to add more samples, pls. read on:
1. Create new sub directory and give it an appropriate name for your new sample. Pls. follow the layout of the existing sample if you are not sure how to organize your code.
2. Make sure your sample work as expected before submit PR, and make sure GitHub CI passes after PR is submitted. Pls. refer to the existing sample on how to test the sample.   
3. Pls. provide README.md to explain your samples.
