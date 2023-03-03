# gRPC example

Dubbo 3.0 provides Triple (Dubbo3) and Dubbo2 protocols, which are native protocols of the Dubbo framework. In addition, Dubbo3 also integrates many third-party protocols and incorporates them into Dubbo's programming and service governance system, including gRPC, Thrift, JSON-RPC, Hessian2, REST, etc.

**This example will introduce how to use the gRPC protocol**.

## Run the example:

Start zk and listen on port 127.0.0.1:2181.
If zk is not installed, you can also use docker to directly execute the following commands to start all dependent components running samples: zk(2181), nacos(8848), etcd(2379).

`docker-compose -f {PATH_TO_SAMPLES_PROJECT}/integrate_test/dockercompose/docker-compose.yml up -d`

### Run via command line

- Server

`cd rpc/grpc/go-server/cmd` # enter the warehouse directory

`export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"`# Set configuration file environment variable

`go run .` # Start the service

- client

`cd rpc/grpc/go-client/cmd` # enter the warehouse directory

`export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"`# Set configuration file environment variable

`go run .` # Start the client to initiate a call

### The call is successful

After the client calls, you can see that the client prints the following information, and the call is successful:

`[XXXX-XX-XX/XX:XX:XX main.main: client.go: 55] client response result: {this is message from reply {} [] 0}`