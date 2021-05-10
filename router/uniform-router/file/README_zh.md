# 使用 Dubbbogo 统一路由规则 示例

## 1. 路由规则文件注解

路由规则只针对客户端，对于服务端，只需要在服务提供时打好特定的参数标签即可。

### 1.1 virtual-service
go-client/conf/virtual_service.yml

```yaml
apiVersion: service.dubbo.apache.org/v1alpha1
kind: VirtualService
metadata: {name: demo-route}
spec:
  dubbo:
    # 使用正则表达式匹配service名，只有个满足该service名的请求才能路由。
    # 就此例子来说，不满足service名的请求会直接找不到provider
    - services:
        - { regex: org.apache.dubbo.UserProvider* }
    - routedetail:
        # 匹配规则，如果（sourceLabel）客户端url满足存在参数 `trafficLabel: xxx` 的才能匹配成功
        - match:
            - sourceLabels: {trafficLabel: xxx}
          name: other-condition
          route: # 一旦匹配上述match规则，将选择dest_rule.yml里定义的名为other的子集
            - destination: {host: demo, subset: other}
        - name: center-match
          # 没有match，兜底逻辑，一定会被匹配到。
          route: # 将选择dest_rule.yml里定义的名为center的子集
            - destination: {host: demo, subset: center}

  hosts: [demo]  # 匹配dest_rule.yml里面的host

```

### 1.2 destination-rule
go-client/conf/dest_rule.yml

```yaml
apiVersion: service.dubbo.apache.org/v1alpha1
kind: DestinationRule
metadata: { name: demo-route }
spec:
  host: demo
  subsets:
    - name: all
      labels: { ut: CENTER } # 选中：服务端url存在 `ut:CENTER` 的键值参数的所有实例作为子集
    - name: center
      labels: { ut: other } # 选中：服务端url存在 `ut:other` 的键值参数的所有实例作为子集
    - name: other # 无条件，选择所有实例
```

## 2. client、server 路由参数设置
- client 端
在本例子中，go-client/conf/client.yml 可看到如下注释
```yaml
# reference config
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    params: { trafficLabel: xxx }
# If this line is comment, the all requests would be send to server, else the request would
# be sent to both server and server2
```
可见，params 对于客户端url参数的定义，一旦增加这个参数，根据上述配置的路由规则，就会命中上述 destination 名为 other 的子集，从而对所有 provider 采用负载均衡策略发起调用。\
而如果注释掉这行参数，会将请求路由至 center 子集，针对单一的 server 发起调用。

在环境变量中配置路由规则文件的路径

```shell
export CONF_VIRTUAL_SERVICE_FILE_PATH=xxx/virtual_service.yml
export CONF_DEST_RULE_FILE_PATH=xxx/dest_rule.yml
```

- server 端
```yaml
# service config
services:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    params: { ut: other } # 服务端标签，destination rule 根据此标签选择子集对应的所有实例
    methods:
      - name: "GetUser"
        retries: 1
        loadbalance: "random"
```

## 3. 运行方法

直接使用goland运行本示例

router/router-server\
router/router-server2\
router/router-client


运行后，可观测到所有客户端流量都路由至 router-server，并没有请求路由至 router-server2

