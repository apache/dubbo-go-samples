# Consul 服务发现示例
### 1. 介绍
[Consul](https://www.consul.io/docs/intro)
是一个服务网格解决方案，它提供了一个具有服务发现、配置和分段功能的全功能控制平面。
在通过dubbogo使用consul之前，请保证consul已经在您的环境正确运行，并且和下面的ip:port配置一致。

### 2. Important config
```yaml
registries: # registry config
  "demoServiceDiscovery": 
    protocol: "service-discovery" # choose service-discovery as service registry config protocol
    params:
      service_discovery: "consul_dis"
      name_mapping: "in-memory"
      metadata: "default"

service_discovery:
  consul_dis: # this key sould match config above
    protocol: "consul" # this is real protocol that sd based
    remote_ref: "consul1" # this is the remoting network config of target protocol

remote: 
  consul1: # this key should match config service_discovery's remote_ref
    address: "127.0.0.1:8500" # consul 服务地址
    timeout: "5s" 

# if you want to report metadata by consul, you can add config below:
metadata_report:
  protocol: "consul"
  remote_ref: "consul1"
```
### 3. Import 引入依赖
使用consul作为服务注册中心，请将引入以下依赖：
```go
import(
    _ "github.com/apache/dubbo-go/metadata/report/consul"
    _ "github.com/apache/dubbo-go/registry/servicediscovery"
)
	
```
如果您希望使用consul来报告数据，请引入以下依赖：
```go
import(
    _ "github.com/apache/dubbo-go/metadata/report/consul"
)
```