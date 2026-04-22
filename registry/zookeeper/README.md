# Zookeeper as registry

This example shows dubbo-go's service discovery feature with Zookeeper as registry.

## How to run

### Start Zookeeper server
This example relies on zookeeper as registry, follow the steps below to start a zookeeper server first.

1. Start zookeeper with docker, run `docker run --rm -p 2181:2181 zookeeper` or `make -f $DUBBO_GO_SAMPLES_ROOT_PATH/Makefile docker-up`.
2. [Download and start zookeeper](https://zookeeper.apache.org/releases.html#download) locally on your machine.

### Run server
```shell
$ go run ./go-server/cmd/server.go
```

test rpc server work as expected:
```shell
$ curl \
    --header "Content-Type: application/json" \
    --data '{"name": "Dubbo"}' \
    http://localhost:20000/greet.GreetService/Greet
```

check url address successfully registered into zookeeper:
```shell
# enter zookeeper bin directory, for example '$HOST_PATH/apache-zookeeper-3.5.9-bin/bin'
$ ./zkCli.sh
[zk: localhost:2181(CONNECTED) 0] ls /services/dubbo_registry_zookeeper_server
[30.221.147.198:20000]
```

### Run client
```shell
$ go run ./go-client/cmd/client.go
hello world
```