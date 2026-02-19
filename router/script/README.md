# Script Router

This example demonstrates how to use the script router feature in dubbo-go.

English | [中文](README_CN.md)

## Prerequisites

- Docker and Docker Compose environment to run the Nacos registry.
- Go version 1.23+.
- Nacos version 2.x+.

## How to Run

### Start Nacos Registry

Follow this instruction to [install and start Nacos server](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

### Run Server (Provider)

In this example, you will run two servers providing services on ports 20000 and 20001 respectively.

```Shell
$ go run ./go-server/cmd/server.go              # port 20000
$ go run ./go-node2-server/cmd/server_node2.go  # port 20001
```

### Run Client (Consumer)

In this example, the client will continuously call the Greet method in an infinite loop. You need to:

- Start the client and observe the load balancing behavior during calls.

- Set up the script router configuration in the Nacos registry, then observe the client's calls again.

```Shell
$ go run ./go-client/cmd/client.go
```

### Nacos Configuration

Create a new configuration with the `Data ID` **script-server.script-router** and set the format to `yaml`.

Set the Group to `DEFAULT_GROUP`.

> Note: The naming convention in Nacos is {application.name}.{router_type}.

```yaml
scope: "application"
key: "condition-server"
enabled: true
type: "javascript"
script: |
  (function(invokers, invocation, context) {
    if (!invokers || invokers.length === 0) return [];
    return invokers.filter(function(invoker) {
      var url = invoker.GetURL();
      return url && url.Port === "20000";
    });
  })(invokers, invocation, context);
```

## Expected Results

- When the client is started without the script router configuration in the Nacos config center, the client will alternate
  calls between the two servers.

- After starting the client and configuring the script router in Nacos, the client will only call one specific server (
  port 20000).