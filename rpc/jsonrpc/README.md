# JsonRPC Example

## Backend
Dubbo3 provides Triple(Dubbo3), Dubbo2 protocols, which are native to the Dubbo framework.
In addition, Dubbo3 also integrates a number of third-party protocols into Dubbo's programming and service governance architecture,
Including gRPC, Thrift, **JsonRPC**, Hessian2, and REST. The following describes the **JsonRPC** protocol example.

## Start

- Start the registry
- Start go-server and go-client, practice with **JsonRPC**  
- Start java-server and java-client, practice with **JsonRPC**

### Start the registry

Start the registry by docker-compose:

```shell
docker-compose -f go-server/docker/docker-compose.yml up -d
```

Stop the registry

```shell
docker-compose -f go-server/docker/docker-compose.yml dowm
```

### Start Go Server and Client

Note: Goland users can directly use the boot mode configured by '.run ', refer to [HOWTO.md](../HOWTO_zh.md)

Start go-server：

Config the configuration file of **Dubbogo**（[server/dubbogo.yml](go-server/conf/dubbogo.yml)）：

```shell
DUBBO_GO_CONFIG_PATH=${$PROJECT_DIR$}/dubbo-go-samples/rpc/jsonrpc/go-server/conf/dubbogo.yml
```

Start go-client：

Config the configuration file of **Dubbogo**（[client/dubbogo.yml](go-client/conf/dubbogo.yml)）：

```shell
DUBBO_GO_CONFIG_PATH=${$PROJECT_DIR$}/dubbo-go-samples/rpc/jsonrpc/go-client/conf/dubbogo.yml
```

### Start Java Server and Client

Start java-server：

run [build.sh](java-server/build.sh) ，Maven environment required

```shell
bash build.sh
```

Start java-client：

run [build.sh](java-client/build.sh)，Maven environment required

```shell
bash build.sh
```


