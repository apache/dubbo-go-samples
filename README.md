# Dubbo Golang Examples

![CI](https://github.com/apache/dubbo-go-samples/workflows/CI/badge.svg)

## What It Contains

* async: An async example. dubbo-go supports client to call server asynchronously.
* attachment: An attachment example, to show how to use attachment to pass user data from the client to the server.  
* configcenter: Some examples of different config center. There are three -- zookeeper, apollo and nacos at present.
* direct: A direct example. This feature make start of dubbo-go get easy.
* filter: Some examples of different filter. Including custom_filter and tpslimit
* general: A general example. It had validated zookeeper registry and different parameter lists of service.
* generic: A generic example. It show how to use generic feature of dubbo-go.
* helloworld: A simplest example. It contain 'go-client', 'go-server', 'java-server' of dubbo protocol. 
* multi_registry: An example of multiple registries.
  And it has a comprehensive testing with dubbo/jsonrpc/grpc/rest protocol. You can refer to it to create your first complete dubbo-go project.
* registry: Some examples of different registry. There are kubernetes, nacos and etcd at present. **Note**: When use the different registry, you need update config file, but also must import the registry package. see the etcd `README`
* router: Some router examples. Now, a condition router example is existing. 
* seata: Transaction system examples by seata.
* shop: Shop sample, make consumer and provider run in a go program.
* tracing: Some tracing examples. We have tracing support of dubbo/grpc/jsonrpc protocol at present. 

## How To Run

Pls. refer [How To Run](HOWTO.md) for the instructions.

## How to contribute

If you want to add some samples, we hope that you can do this:
1. Adding samples in appropriate directory. If you dont' know which directory you should put your samples into, you can get some advices from dubbo-go community.
2. You must run the samples locally and there must be no any error.
3. If your samples have some third party dependency, including another framework, we hope that you can provide some docs, script is better.
