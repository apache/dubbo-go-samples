# Docker example
### 1. How to run
This example shows developers how to build your server using dubbo-go framework to docker. 

run :
```shell script
sudo sh docker_run.sh
``` 

Then you can build and pack your server to running docker.

### docker file:
```dockerfile
FROM golang:latest as build
ENV GO111MODULE on
ENV CONF_PROVIDER_FILE_PATH  /release/conf/server.yml
ADD ./go-server/dist/linux_amd64/ /
ENTRYPOINT exec /release/go-server
```

### 2. Attention

in docker_run.sh
```shell script
P_DIR=$(pwd)/go-server
make GOOS="linux" PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../build/Makefile build
docker build --no-cache -t dubbogo-docker-sample .
docker run --name zkserver -p 2181:2181 --restart always -d zookeeper:3.4.9
docker run -e DUBBO_IP_TO_REGISTRY=127.0.0.1  -p 20000:20000  --link zkserver:zkserver dubbogo-docker-sample
```

- first start zookeeper as registry
- second we should set ip_to_registry environment.  \
Docker bridge network mode need to specify a registered host IP for external network communication. Dubbo provides two pairs of system attributes in the startup phase, which are used to set the IP and port addresses of external communication. 

### 3. Test
in go-server/tests/integration/

run:
```shell script
go test -tags integration -v .
```          

You can get lines at the end of output:
```
--- PASS: TestGetUser (0.00s)
PASS
ok      github.com/apache/dubbo-go-samples/docker/go-server/tests/integration   3.612s
```