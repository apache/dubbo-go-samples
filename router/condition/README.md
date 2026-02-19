# Condition router

This example shows how to use dubbo-go's condition router.

English | [中文](README_CN.md)

## Prerequisites

- Docker and docker compose to run Nacos registry.
- Go Version 1.23+.
- Nacos Version 2.x+.

## How to run

### Run Nacos registry

Follow this instruction
to [install and run Nacos](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

### Run server (Provider)

In this example, you will run two servers on ports 20000 and 20001 respectively.

```shell
$ go run ./go-server/cmd/server.go              # port 20000
$ go run ./go-node2-server/cmd/server_node2.go  # port 20001
```

### Run client (Consumer)

In this example, the client will keep calling Greet method in an infinity loop. You need to:

- Start the client and observe the load balancing behavior during calls.

- Set up the condition router configuration in the Nacos registry, then observe the client's calls again.

```shell
$ go run ./go-client/cmd/client.go
```

### Configuration of Nacos

Create a new configuration with the `Data ID` **condition-server.condition-router** and set the format to `yaml`.

Set the Group to `DEFAULT_GROUP`.

> Note: The naming convention in Nacos is {application.name}.{router_type}.

```yaml
configVersion: V3.3.2
scope: "application"
key: "condition-server"
priority: 1
force: true
enabled: true
conditions:
  - from:
      match: "application = condition-client"
    to:
      - match: "port = 20001"
```

## Expected result

- When the client starts without the condition router configuration in Nacos,
  it will alternate calls between the two server endpoints.
- When the client starts with the script router configured in Nacos,
  it will route traffic to only one server endpoint.
