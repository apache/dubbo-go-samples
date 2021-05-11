# Using Dubbbogo unified routing rules example

## 1. Notes on routing rules file

The routing rules are only for the client. For the server, you only need to put a specific parameter label when the service is provided.

### 1.1 virtual-service
go-client/conf/virtual_service.yml

```yaml
apiVersion: service.dubbo.apache.org/v1alpha1
kind: VirtualService
metadata: {name: demo-route}
spec:
  dubbo:
    # Use regular expressions to match the service name, and only a request that meets the service name can be routed.
    # For this example, a request that does not satisfy the service name will directly not find the provider
    -services:
        -{regex: org.apache.dubbo.UserProvider*}
    -routedetail:
        # Matching rules, if the (sourceLabel) client url satisfies the existence of the parameter `trafficLabel: xxx`, the match can be successful
        -match:
            -sourceLabels: {trafficLabel: xxx}
          name: other-condition
          route: # Once the above match rules are matched, the subset named other defined in dest_rule.yml will be selected
            -destination: {host: demo, subset: other}
        -name: center-match
          # There is no match, and the logic will be matched.
          route: # The subset named center defined in dest_rule.yml will be selected
            -destination: {host: demo, subset: center}

  hosts: [demo] # match the host in dest_rule.yml

```

### 1.2 destination-rule
go-client/conf/dest_rule.yml

```yaml
apiVersion: service.dubbo.apache.org/v1alpha1
kind: DestinationRule
metadata: {name: demo-route}
spec:
  host: demo
  subsets:
    -name: all
      labels: {ut: CENTER} # Selected: all instances of the key-value parameter of `ut:CENTER` in the server url as a subset
    -name: center
      labels: {ut: other} # Checked: all instances of the key-value parameter of `ut:other` in the server url as a subset
    -name: other # Unconditionally, select all instances
```

## 2. Client and server routing parameter settings
-client side
In this example, go-client/conf/client.yml can see the following comments
```yaml
# reference config
references:
  "UserProvider":
    registry: "demoZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    params: {trafficLabel: xxx}
# If this line is comment, the all requests would be send to server, else the request would
# be sent to both server and server2
```
It can be seen that params defines the client url parameter. Once this parameter is added, according to the routing rules configured above, it will hit the subset of the destination named other, so that the load balancing strategy is used to initiate calls to all providers. \
If you comment out this line of parameters, the request will be routed to the center subset, and a call will be initiated for a single server.

-server side
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
    params: {ut: other} # server label, the destination rule selects all instances corresponding to the subset according to this label
    methods:
      -name: "GetUser"
        retries: 1
        loadbalance: "random"
```

## 3. How to run

Run this example using goland

router/router-server\
router/router-server2\
router/router-client


After running, it can be observed that all client requests are routed to router-server, and no request is routed to router-server2 