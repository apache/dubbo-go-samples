# Condition router
This example shows how to use dubbo-go's condition router.

English | [中文](README_CN.md)

## Prerequisites

- Docker and docker compose to run Nacos registry.
- Go 1.23+.
- Nacos 2.x+.

## How to run

### Run Nacos registry
Follow this instruction to [install and run Nacos](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

### Run server (Provider)
In this example, you will run two servers with one running on port 20000 and another 20001.

```shell
$ go run ./go-server/cmd/server-base/server.go        # port 20000
$ go run ./go-server/cmd/server-copy/server_copy.go   # port 20001
```

### Run client (Consumer)
In this example, the client will keep calling Greet method in an infinity loop, and you need to:

- run client, and observe its calling result (with load balance).
- configure `condition router` configuration in Nacos, and observe the result of calling.

```shell
$ go run ./go-client/cmd/client.go
```

### Configuration of Nacos

Create a new configuration with `Data ID` in `condition-server.condition-router` and in `yaml` extension.

you need to set its group in `DEFAULT_GROUP`.

> Tips: In Nacos, your router configuration name should in such formation {application.name}.{router_type}

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

- When the client starts but the configuration for condition router is not set in the Nacos configuration center, 
  the client will call back and forth between two server endpoints.
- When the client starts with configuration for condition router in the Nacos, the client will only call one server endpoint.
