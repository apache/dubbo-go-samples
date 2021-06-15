In order to make the client work, a conf/server.yml shown below is required. Need further investigation:

```yaml
# dubbo server yaml configure file

# application config
application:
  organization: "dubbo.io"
  name: "UserInfoServer"
  module: "dubbo-go user-info server"
  version: "0.0.1"
  environment: "dev"

# registry config
registries:
  "demoServiceDiscovery":
    protocol: "service-discovery"
    params:
      service_discovery: "nacos1"
      name_mapping: "dynamic"
      metadata: "default"

remote:
  nacos:
    address: "127.0.0.1:8848"
    timeout: "5s"

metadata_report:
  protocol: "nacos"
  remote_ref: "nacos"

service_discovery:
  nacos1:
    protocol: "nacos"
    remote_ref: "nacos"

config_center:
  protocol: "nacos"
  address: "127.0.0.1:8848"
```