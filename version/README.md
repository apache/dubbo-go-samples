# Multiple Versions

When an interface to achieve an incompatible upgrade, you can use the version number transition. Different versions of the services do not reference each other.

You can follow the steps below for version migration:

0. In the low pressure period, upgrade to half of the provider to the new version
0. Then upgrade all consumers to the new version
0. Then upgrade the remaining half providers to the new version

## Introduced

```markdown
.
├── README.md
├── README_zh.md
├── docker
├── go-api
├── go-client
├── go-server-v1
└── go-server-v2

```

- go-client ：The Service Consumer
- go-server-v1 ：The Service Provider(Old Version)
- go-server-v2 ：The Service Provider(New Version)
- go-api ：Public

### Config Instructions

Consume the specified service by changing the Provider and Consumer version numbers

Consumer:

```yaml
references:
  "UserProvider":
    version: "1.0" # 2.0
```

Provider:

```yaml
services:
  "UserProvider":
    version: "1.0" # 2.0
```

#### Provider Configuration

Old version of the service provider configuration:

```yaml
# service config
services:
  "UserProvider":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    version: "1.0"
    methods:
    - name: "GetUser"
      retries: 1
      loadbalance: "random"
```

New version of the service provider configuration:

```yaml
# service config
services:
  "UserProvider":
    registry: "demoZk"
    protocol : "dubbo"
    interface : "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    version: "2.0"
    methods:
    - name: "GetUser"
      retries: 1
      loadbalance: "random"
```

#### Consumer Configuration

Old version of the service consumer configuration:

```yaml
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    version: "1.0"
    methods:
      - name: "GetUser"
        retries: 3
```

New version of the service consumer configuration:

```yaml
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    version: "2.0"
    methods:
      - name: "GetUser"
        retries: 3
```

#### Result

Consumer log:

```markdown
response result: {A001 Alex Stocks 18 2021-05-09 18:30:16.957 +0800 CST Provider Version 2.0}
```

Provider log:

```markdown
Server V2 req:[]interface {}{"A001"}[2021-05-09/18:30:16 github.com/apache/dubbo-go-samples/version/go-server-v2/pkg.(*UserProvider).GetUser: user.go: 47] %s
Server V2 rsp:pkg.User{ID:"A001", Name:"Alex Stocks", Age:18, Time:time.Time{wall:0xc01e0c4e39144248, ext:20634030436, loc:(*time.Location)(0x4b80960)}, ServiceInfo:"Provider Version 2.0"}
```

### How To Run

Refer to  [HOWTO.md](../HOWTO_zh.md) under the root directory to run this sample.
