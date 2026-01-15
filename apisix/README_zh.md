# APISIX 整合 Dubbo-Go

[English](README.md) | [中文](README_zh.md)

本示例演示如何使用 Apache APISIX 作为 Dubbo-Go 服务的 API 网关。

## 概述

Apache APISIX 是一个动态、实时、高性能的 API 网关，本示例展示了如何：
- 使用 Dubbo-Go 3.x 的新 API 创建 Triple 协议服务
- 通过 Nacos 进行服务注册与发现
- 使用 APISIX 网关进行 HTTP 到 gRPC/Triple 的协议转换
- 通过 APISIX 路由访问 Dubbo-Go 服务

## 环境要求

- Linux 或 macOS
- Docker 20.10+
- Docker Compose v2.0+
- Go 1.18+ (如需本地开发)

## 架构说明

```
Client (HTTP) -> APISIX Gateway -> Dubbo-Go Service (Triple/gRPC)
                       ↓                    ↓
                     etcd              Nacos Registry
```

## 快速开始

### 1. 创建 Docker 网络

```bash
docker network create default_network
```

### 2. 启动依赖服务

按顺序启动 etcd、MySQL、Nacos 和 APISIX：

```bash
# 启动 etcd (APISIX 的配置中心)
cd ./deploy/etcd-compose
docker compose up -d

# 启动 MySQL (Nacos 的数据库)
cd ../mysql5.7-compose
docker compose up -d

# 等待 MySQL 启动完成 (约10秒)
sleep 10

# 启动 Nacos
cd ../nacos2.0.3-compose
docker compose up -d

# 等待 Nacos 启动完成 (约20秒)
sleep 20

# 启动 APISIX (API 网关)
cd ../apisix-compose
docker compose up -d

cd ../../
```

### 3. 构建并启动 Dubbo-Go 服务

```bash
# 构建镜像
./build.sh

# 获取 Nacos 容器 IP
NACOS_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' nacos203-compose-nacos-1)
echo $NACOS_IP

# 启动 Dubbo-Go 服务
docker run -d \
  --name dubbo-go-apisix-server \
  --network default_network \
  -e NACOS_ADDR=$NACOS_IP:8848 \
  dubbo-go-apisix-server:latest
```

### 4. 配置 APISIX

#### 4.1 配置 Protocol Buffer

```bash
curl -X PUT 'http://127.0.0.1:80/apisix/admin/proto/1?api_key=edd1c9f034335fi23f87ad84b625c8f1' \
-H 'Content-Type: application/json' \
-d '{
    "content": "syntax = \"proto3\";\npackage helloworld;\n\noption go_package = \"github.com/apache/dubbo-go-samples/apisix/proto;greet\";\n\nservice Greeter {\n  rpc SayHello (HelloRequest) returns (User) {}\n  rpc SayHelloStream (stream HelloRequest) returns (stream User) {}\n}\n\nmessage HelloRequest {\n  string name = 1;\n}\n\nmessage User {\n  string name = 1;\n  string id = 2;\n  int32 age = 3;\n}"
}'
```

#### 4.2 配置路由

```bash
curl -X PUT 'http://127.0.0.1:80/apisix/admin/routes/1?api_key=edd1c9f034335fi23f87ad84b625c8f1' \
-H 'Content-Type: application/json' \
-d '{
    "uri": "/helloworld",
    "name": "helloworld",
    "methods": ["GET", "POST"],
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
        "service_name": "dubbo_apisix_server"
    },
    "status": 1
}'
```

### 5. 测试

```bash
curl 'http://127.0.0.1:80/helloworld?name=World'
```

预期输出：

```json
{
    "age": 21,
    "id": "12345",
    "name": "Hello World"
}
```

## 验证服务注册

访问 Nacos 控制台查看服务注册情况：

```
http://localhost:8848/nacos
用户名: nacos
密码: nacos
```

在服务列表中应该能看到 `providers:helloworld.Greeter::` 服务。

## 参考文档

- [Apache APISIX 文档](https://apisix.apache.org/zh/docs/apisix/getting-started)
- [Dubbo-Go 官方文档](https://dubbogo.apache.org/)
- [Nacos 官方文档](https://nacos.io/)
