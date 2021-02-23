# Docker 示例
### 1. 如何使用
这个示例展示了开发者如何将dubbogo框架开发的应用打包为docker，并提供测试方案。

当前目录下终端运行

```shell script
sudo sh docker_run.sh
```

你可以打包构建你的服务，并且以docker方式运行

### Dockerfile:
```dockerfile
FROM golang:latest as build
ENV GO111MODULE on
ENV CONF_PROVIDER_FILE_PATH  /release/conf/server.yml
ADD ./go-server/dist/linux_amd64/ /
ENTRYPOINT exec /release/go-server
```

### 2.  重要过程

在 docker_run.sh 脚本中
```shell script
P_DIR=$(pwd)/go-server
make GOOS="linux" PROJECT_DIR=$P_DIR PROJECT_NAME=$(basename $P_DIR) BASE_DIR=$P_DIR/dist -f ../build/Makefile build
docker build --no-cache -t dubbogo-docker-sample .
docker run --name zkserver -p 2181:2181 --restart always -d zookeeper:3.4.9
docker run -e DUBBO_IP_TO_REGISTRY=127.0.0.1  -p 20000:20000  --link zkserver:zkserver dubbogo-docker-sample
```

- 首先开启了zookeeper docker，用于注册服务

- 之后需要设置t ip_to_registry 环境变量.  \

  docker桥接网络模型需要明确指定主机ip用于网络通信。Dubbo框架在运行参数中提供了两个系统属性，能够指定Ip和port进行注册，然后其他服务通过注册好的信息访问当前服务。

### 3. 测试
在 go-server/tests/integration/ 文件夹下

运行:
```shell script
go test -tags integration -v .
```

你可以在测试日志的最后，看到测试成功的信息：
```
--- PASS: TestGetUser (0.00s)
PASS
ok      github.com/apache/dubbo-go-samples/docker/go-server/tests/integration   3.612s
```
