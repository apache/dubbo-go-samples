# Generic Invocation

Generic invocation ensures the RPC could be passed properly even if one of clients has no information about interface, because generic invocation converts the POJO to a generic type, like dictionary, string, etc. It is often used for testing and gateways. Please visit our documentation for more details.

## Getting Started

The samples of generic invocation are parted by the way of generalization：

- default：uses MapGeneralizer which converts POJOs to maps

Each sample contains 4 subfolders：

- go-server：Dubbo-Go server sample
- go-client：Dubbo-Go client sample
- java-client：Dubbo server sample
- java-server：Dubbo client sample

Providing java samples is convenient to test interoperability between Dubbo and Dubbo-Go.

### Registry

This sample uses ZooKeeper as the registry. In fact, etcd and Nacos are supported as well. Executing the following command, a ZooKeeper instance will be launched. Please note that docker and docker-compose **SHOULD** be installed before.

```shell
cd ./default/go-server/docker \
  && docker-compose up -d
```
### Server

There are two ways to launch a Dubbo-Go server: using GoLand or using command line tool.

Using GoLand. Please select `v3config-generic/generic-default-go-server` from Configurations at top-right corner, and then click Run button.

Using command line tool. The `$ProjectRootDir` is the root directory of the dubbo-go-samples project.

```shell
cd $ProjectRootDir/generic/default/go-server/cmd \
  && go run server.go
```

### Client

There are two ways to launch a Dubbo-Go client: using GoLand or using command line tool.

Using GoLand. Please select `v3config-generic/generic-default-go-client` from Configurations at top-right corner, and then click Run button.

Using command line tool. The `$ProjectRootDir` is the root directory of the dubbo-go-samples project.

```shell
cd $ProjectRootDir/generic/default/go-client/cmd \
  && go run client.go
```
