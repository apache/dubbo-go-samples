# APISIX 整合 Dubbo Go



**[Demo](https://github.com/limerence-code/apisix-dubbo-go.git)**

## 环境准备

1. Linux
2. docker
3. docker-compose

本文以 **Ubuntu 22.04 LTS** , **docker 20.10.14**,  **docker-compose v2.2.2** 为例

## 依赖

1. apisix  
2. apisix-dashboard 
3. etcd  
4. helloword  
5. nacos  

### apisix 

APISIX 服务端

### apisix-dashboard (可选)

APISIX 控制台，提供可视化控制

### etcd

APISIX 的注册中心

### helloworld

Dubbo Go 的测试 Demo

### nacos

用于注册 Dubbo Go 服务，供 APISIX 网关调用

## 启动

### 创建 docker network

```shell
docker network create default_network
```

创建 default_network ，服务指定该网络；方便服务之间进行通讯

### 依次启动服务

按顺序启动 **etcd** 、**apisix** 、  **nacos**  、 **helloworld** 命令 **docker-compose up --build -d**

如果需要通过控制台进行协议路由配置则可以启动 **apisix-dashboard** 本文介绍的是通过 HTTP 直接控制，因此无需启动

**PS: 启动 helloworld 服务时，需要提前查询 nacos 对应 default_network 中的 ip ，然后将 main.go 中 nacosConfig.Address 修改成对应的 nacos 地址**

```shell
docker inspect --format='{{json .NetworkSettings.Networks}}'  nacos
```

helloworld 启动成功后，在 nacos 服务列表可以查看

## 配置

### 协议配置

```apl
curl --location --request PUT 'http://127.0.0.1:80/apisix/admin/proto/1?api_key=edd1c9f034335fi23f87ad84b625c8f1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "content": "syntax = \"proto3\";\npackage helloworld;\n\noption go_package = \"./;helloworld\";\n\n// The greeting service definition.\nservice Greeter {\n  // Sends a greeting\n  rpc SayHello (HelloRequest) returns (User) {}\n  // Sends a greeting via stream\n  rpc SayHelloStream (stream HelloRequest) returns (stream User) {}\n}\n\n// The request message containing the user'\''s name.\nmessage HelloRequest {\n  string name = 1;\n}\n\n// The response message containing the greetings\nmessage User {\n  string name = 1;\n  string id = 2;\n  int32 age = 3;\n}"
}'
```

其中 content 内容就是 helloworld.proto 内容， api_key 在 apisix_conf 下面即可找到

配置了协议 id 为1的协议，下面会用到

### 路由转发

```apl
curl --location --request PUT 'http://127.0.0.1:80/apisix/admin/routes/1?api_key=edd1c9f034335fi23f87ad84b625c8f1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "uri": "/helloworld",
    "name": "helloworld",
    "methods": [
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "PATCH",
        "HEAD",
        "OPTIONS",
        "CONNECT",
        "TRACE"
    ],
    "plugins": {
        "grpc-transcode": {
            "method": "SayHello",
            "proto_id": "1",
            "service": "helloworld.Greeter"
        }
    },
    "upstream": {
        "type": "roundrobin",
        "scheme": "grpc",
        "discovery_type": "nacos",
        "pass_host": "pass",
        "service_name": "providers:helloworld.Greeter::"
    },
    "status": 1
}'
```

以上配置表示通过 /helloworld ，可以路由到 helloworld.Greeter 中的 SayHello 方法

详细配置可查看 [APISIX](https://apisix.apache.org/zh/docs/apisix/getting-started)

## 访问

```api
curl --location --request GET 'http://127.0.0.1:80/helloworld?api_key=edd1c9f034335f136f87ad84b625c8f1'
```

输出

```json
{
    "age": 21,
    "id": "12345",
    "name": "Hello "
}
```

