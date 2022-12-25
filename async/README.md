# Async RPC Dubbo for Dubbo-go

For api definition and go client and server startup, please refer to [dubbo-go 3.0 quickstart](https://dubbogo.github.io/zh-cn/docs/user/quickstart/3.0/quickstart.html)

## Instructions
### 1. Start zookeeper
Execute the command `docker run --rm -p 2181:2181 zookeeper` or `make -f $DUBBO_GO_SAMPLES_ROOT_PATH/build/Makefile docker-up`.
   If you choose the second way, please ensure that you have set the environment $DUBBO_GO_SAMPLES_ROOT_PATH.
   
### 2. Start the server

Use goland to start async/go-server/cmd/server.go

### 3. Start the client

Use goland to start async/go-client/cmd/client.go