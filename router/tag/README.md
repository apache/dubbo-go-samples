# Tag router
This example shows how to use dubbo-go tag router ability.

## Prerequisites

- Docker and Docker Compose for running Nacos registry.
- Go 1.23+ for Dubbo-Go examples.

## How to run

### Start Nacos server
Follow this instruction to [install and start Nacos server](https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/).

### Run server
In this example, you need to run two different server with one with tag and another without.

```shell
$ go run ./go-server/cmd/server.go          # non-tagged server
$ go run ./go-tag-server/cmd/server_tag.go  # tagged server
```

### Run client

```shell
$ go run ./go-client/cmd/client.go
```

## Expected outputs

- ✔ invoke successfully : receive: tag with force, response from: server-with-tag
- ❌ invoke failed: Failed to invoke the method Greet.
- ✔ invoke successfully : receive: tag with no-force, response from: server-without-tag
- ✔ invoke successfully : receive: non-tag, response from: server-without-tag
