# Consul Service Discovery Example
### 1. Introduction
[Consul](https://www.consul.io/docs/intro) is a service mesh solution providing a full featured control plane with service discovery, configuration, and segmentation functionality.\
Dubbogo can use consul as service discovery middleware and metadata reporter. \
Before using consul, pls make sure consul server is running in your environment and match the config below.\

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
    address: "127.0.0.1:8500" # consul address
    timeout: "5s" 

# if you want to report metadata by consul, you can add config below:
metadata_report:
  protocol: "consul"
  remote_ref: "consul1"
```
### 3. Import block
To use consul as service discovery, make sure you add import in you server and client:
```go
import(
    _ "github.com/apache/dubbo-go/metadata/report/consul"
    _ "github.com/apache/dubbo-go/registry/servicediscovery"
)
	
```

if you want to use consul to report, make sure you add import:
```go
import(
    _ "github.com/apache/dubbo-go/metadata/report/consul"
)
```