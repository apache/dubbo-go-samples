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
      service_discovery: "file1"
      name_mapping: "in-memory"
      metadata: "default"

service_discovery:
  file1:
    protocol: "file"
```