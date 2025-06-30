# Dubbo Go & Polaris Registry Example

English | [中文](README-zh.md)

## Using the service registration discovery function

Quickly experience Polaris' service registration and service discovery capabilities in dubbogo

## Polaris server installation

[Polaris Server Standalone Version Installation Documentation](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E5%8D%95%E6%9C%BA%E7%89%88%E5%AE%89%E8%A3%85/)

[Polaris Server Cluster Version Installation Documentation](https://polarismesh.cn/docs/%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97/%E6%9C%8D%E5%8A%A1%E7%AB%AF%E5%AE%89%E8%A3%85/%E9%9B%86%E7%BE%A4%E7%89%88%E5%AE%89%E8%A3%85/)

## how to use

### dubbogo.yaml configuration file

Currently PolarisMesh has implemented the registration discovery extension point of dubbogo, so you only need to adjust the registries configuration item in your dubbogo.yaml file, and add the registry configuration of polaris as the protocol. You can refer to the following example.

````yaml
dubbo:
  registries:
    polarisMesh:
      protocol: polaris
      address: ${Polaris server IP}:8091
      namespace: ${Polaris namespace information}
      token: ${Polaris resource authentication token} # If the Polaris server has enabled authentication for the client, you need to configure this parameter
````

### Running the service provider

Enter the cmd directory of go-server and execute the following command

````
 export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
 go run .
````

When you see the following log, it means that the server side started successfully

````log
INFO dubbo/dubbo_protocol.go:84 [DUBBO Protocol] Export service:
````


### Run the service caller

Enter the cmd directory of go-client and execute the following command


````
 export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
 go run .
````

When you see the following log, it means that go-client successfully discovered go-server and made an RPC call

````log
INFO cmd/main.go:75 response: &{A001 Alex Stocks 18 2022-11-19 12:52:38.092 +0800 CST}
````