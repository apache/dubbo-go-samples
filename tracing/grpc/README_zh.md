### 1. 利用docker运行jaeger

[jaeger-getting-started](https://www.jaegertracing.io/docs/1.17/getting-started/)

使用 [all-in-one](https://hub.docker.com/r/jaegertracing/all-in-one) 镜像

### 2. 启动zookeeper

### 3. 启动服务端

设置Jaeger环境变量
```
CONF_PROVIDER_FILE_PATH=xxxxxxxxxxxxx
JAEGER_AGENT_PORT=32769
JAEGER_AGENT_HOST=localhost
JAEGER_SERVICE_NAME=GrpcServer
JAEGER_SAMPLER_PARAM=1
```

具体细节详见：[jaeger-environment-variables](https://github.com/jaegertracing/jaeger-client-go#environment-variables)

另：
- JAEGER_SAMPLER_PARAM必须要设置成```1```，意味着100%的请求会被用于取样，如果是```0.9```就是90%会被取样
- JAEGER_AGENT_PORT=32769，32769是docker发布的端口，这个端口会被映射到6831，golang客户端需要使用6831来发送追踪数据。

然后运行服务端。

### 4. 运行客户端

设置Jaeger环境变量
```
CONF_PROVIDER_FILE_PATH=xxxxxxxxxxxxx
JAEGER_AGENT_PORT=32769
JAEGER_AGENT_HOST=localhost
JAEGER_SERVICE_NAME=GrpcServer
JAEGER_SAMPLER_PARAM=1
```

然后启动客户端，见[README](https://github.com/dubbogo/dubbo-samples/blob/master/golang/README.md)。

### 5. 在JAEGER-UI中查看追踪数据

打开[http://localhost:32768/search](http://localhost:32768/search)