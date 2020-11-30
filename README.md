# Dubbo Golang Examples

## What It Contains

* helloworld: A simplest example. It contain 'go-client', 'go-server', 'java-server' of dubbo protocol. 
* direct: A direct example. This feature make start of dubbo-go get easy. 
* async: An async example. dubbo-go supports client to call server asynchronously. 
* general: A general example. It had validated zookeeper registry and different parameter lists of service. 
  And it has a comprehensive testing with dubbo/jsonrpc/grpc/rest protocol. You can refer to it to create your first complete dubbo-go project.
* generic: A generic example. It show how to use generic feature of dubbo-go.
* configcenter: Some examples of different config center. There are three -- zookeeper, apollo and nacos at present.
* multi_registry: An example of multiple registries.
* registry: Some examples of different registry. There are kubernetes, nacos and etcd at present. **Note**: When use the different registry, you need update config file, but also must import the registry package. see the etcd `README`
* filter: Some examples of different filter. Including custom_filter and tpslimit
* router: Some router examples. Now, a condition router example is existing. 
* seata: Transaction system examples by seata.
* shop: Shop sample, make consumer and provider run in a go program.
* tracing: Some tracing examples. We have tracing support of dubbo/grpc/jsonrpc protocol at present. 

## How To Run

> Take `helloworld` as an example

#### 1. Setup Zookeeper Server

A [zookeeper](https://zookeeper.apache.org/releases.html) server is required to run most of the samples in this repository. It can either start without docker environment like this:

```bash
zkServer start
```

or start when docker environment presents:

```bash
docker run --name zookeeper -p2181:2181 -d zookeeper
```

This samples repository provides an even more convenient way to start zookeeper:

```bash
cd helloworld/go-server
make -f ../../build/Makefile docker-up
```

Once the following messages outputs, the zookeeper server is ready.

```bash
>  Starting dependency services with docker/docker-compose.yml
Creating network "docker_default" with the default driver
Creating docker_zookeeper_1 ... done
```

To shutdown it, simple run

```bash
make -f ~/github/apache/dubbo-go-samples/build/Makefile docker-down
```

#### 2. Start go-server

Use the following commands to start `go-server`.

```bash
cd helloworld/go-server
make -f ../../build/Makefile start
```

Once the following messages outputs, the server is ready.

```bash
>  Buiding application binary: dist/darwin_amd64/release/go-server
>  Starting application go-server, output is redirected to dist/darwin_amd64/release/go-server.log
  >  PID: 86428
```

The output of `go-server` can be found from 'dist/darwin_amd64/release/go-server.log'.

#### 3. Run go-client

Use the following commands to run `go-client`.

```bash
cd helloworld/go-client
make -f ../../build/Makefile run
```

Once the following messages outputs, the `go-client` calls the `go-server` successfully.

```bash
>  Buiding application binary: dist/darwin_amd64/release/go-client
>  Running application go-client, output is redirected to dist/darwin_amd64/release/go-client.log
...
2020-10-27T14:51:37.520+0800    DEBUG   dubbo/dubbo_invoker.go:144      result.Err: <nil>, result.Rest: &{A001 Alex Stocks 18 2020-10-27 14:51:37.52 +0800 CST}
2020-10-27T14:51:37.520+0800    DEBUG   proxy/proxy.go:177      [makeDubboCallProxy] result: &{A001 Alex Stocks 18 2020-10-27 14:51:37.52 +0800 CST}, err: <nil>
response result: &{A001 Alex Stocks 18 2020-10-27 14:51:37.52 +0800 CST}
```

#### 4. Integration Test

dubbo-go-samples is designed to serve the purposes of not only the show cases of how to use apache/dubbo-go but also the integration-test for apache/dubbo-go. To run integration test for `go-server`, run the following commands:

```bash
cd helloworld/go-server
make -f ../../build/Makefile integration
```

Once the following messages outputs, the integration tests pass.

```bash
>  Running integration test for application go-server
...
--- PASS: TestGetUser (0.00s)
PASS
ok      github.com/apache/dubbo-samples/golang/helloworld/go-server/tests/integration   3.603s
```

#### 5. Clean Up

To clean up, run the following command:

```bash
cd helloworld/go-server
make -f ../../build/Makefile clean
make -f ../../build/Makefile docker-down
```

## How to debug with Goland

#### 1. Edit Configurations

![](.images/edit_configuratios.png)

#### 2. Configure `Environment Variable`

* Add `APP_LOG_CONF_FILE`. eg: `/home/xx/dubbogo-samples/helloworld/client/conf/log.yml`
* Add `CONF_CONSUMER_FILE_PATH` eg: `/home/xx/dubbogo-samples/helloworld/client/conf/client.yml`
* Add `CONF_PROVIDER_FILE_PATH` eg: `/home/xx/dubbogo-samples/helloworld/server/conf/server.yml`
![](.images/edit_env.png)
	
Click Apply, then you are all set to run.

## How to contribute

If you want to add some samples, we hope that you can do this:
1. Adding samples in appropriate directory. If you dont' know which directory you should put your samples into, you can get some advices from dubbo-go community.
2. You must run the samples locally and there must be no any error.
3. If your samples have some third party dependency, including another framework, we hope that you can provide some docs, script is better.
