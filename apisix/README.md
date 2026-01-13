# APISIX Integration with Dubbo-Go

This example demonstrates how to use Apache APISIX as an API gateway for Dubbo-Go services.

## Overview

Apache APISIX is a dynamic, real-time, high-performance API gateway. This example shows how to:
- Create Triple protocol services using Dubbo-Go 3.x new API
- Use Nacos for service registration and discovery
- Use APISIX gateway for HTTP to gRPC/Triple protocol conversion
- Access Dubbo-Go services through APISIX routes

## Requirements

- Linux or macOS
- Docker 20.10+
- Docker Compose v2.0+
- Go 1.18+ (for local development)

## Architecture

```
Client (HTTP) -> APISIX Gateway -> Dubbo-Go Service (Triple/gRPC)
                       ↓                    ↓
                     etcd              Nacos Registry
```

## Quick Start

### 1. Create Docker Network

```bash
docker network create default_network
```

### 2. Start Dependencies

Start etcd, MySQL, Nacos, and APISIX in order:

```bash
# Start etcd (APISIX configuration center)
cd ./deploy/etcd-compose
docker compose up -d

# Start MySQL (Nacos database)
cd ../mysql5.7-compose
docker compose up -d

# Wait for MySQL to be ready (~10 seconds)
sleep 10

# Start Nacos
cd ../nacos2.0.3-compose
docker compose up -d

# Wait for Nacos to be ready (~20 seconds)
sleep 20

# Start APISIX (API gateway)
cd ../apisix-compose
docker compose up -d

cd ../../
```

### 3. Build and Start Dubbo-Go Service

```bash
# Build image
./build.sh

# Get Nacos container IP
NACOS_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' nacos)

# Start Dubbo-Go service
docker run -d \
  --name dubbo-go-apisix-server \
  --network default_network \
  -e NACOS_ADDR=$NACOS_IP:8848 \
  dubbo-go-apisix-server:latest
```

### 4. Configure APISIX

#### 4.1 Configure Protocol Buffer

```bash
curl -X PUT 'http://127.0.0.1:80/apisix/admin/proto/1?api_key=edd1c9f034335fi23f87ad84b625c8f1' \
-H 'Content-Type: application/json' \
-d '{
    "content": "syntax = \"proto3\";\npackage helloworld;\n\noption go_package = \"github.com/apache/dubbo-go-samples/apisix/proto;greet\";\n\nservice Greeter {\n  rpc SayHello (HelloRequest) returns (User) {}\n  rpc SayHelloStream (stream HelloRequest) returns (stream User) {}\n}\n\nmessage HelloRequest {\n  string name = 1;\n}\n\nmessage User {\n  string name = 1;\n  string id = 2;\n  int32 age = 3;\n}"
}'
```

#### 4.2 Configure Route

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

### 5. Test

```bash
curl 'http://127.0.0.1:80/helloworld?name=World'
```

Expected output:

```json
{
    "age": 21,
    "id": "12345",
    "name": "Hello World"
}
```

## Verify Service Registration

Access Nacos console to check service registration:

```
http://localhost:8848/nacos
Username: nacos
Password: nacos
```

You should see the `providers:helloworld.Greeter::` service in the service list.

## References

- [Apache APISIX Documentation](https://apisix.apache.org/docs/apisix/getting-started)
- [Dubbo-Go Documentation](https://dubbogo.apache.org/)
- [Nacos Documentation](https://nacos.io/)
