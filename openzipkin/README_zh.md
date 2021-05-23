# Zipkin 在 Dubbo-go 应用示例

# 背景

Zipkin是一个分布式跟踪系统。它有助于收集解决 Dubbo-go 服务中的延迟问题所需的时序数据。
包括该数据的收集和查找。

# 介绍

本示例演示了 Zipkin 在 Dubbo-go 应用程序中的基本用法。

## 目录

```markdown
.
├── README.md
├── README_zh.md
├── docker-compose.yml
├── go-client
├── go-server-a
├── go-server-b
└── prometheus
```

- go-client ：服务消费者
- go-server-a ：服务提供者 A
- go-server-b ：服务提供者 B

## 安装 Zipkin

第一种：

参考 [Zipkin's quick start](https://zipkin.io/pages/quickstart.html) 安装 Zopkin

Follow [Zipkin's quick start](https://zipkin.io/pages/quickstart.html) to install zipkin.

```bash
curl -sSL https://zipkin.io/quickstart.sh | bash -s
```

Zipkin支持多种后端存储，包括Cassandra, ElasticSearch和MySQL。这里我们使用最简单的内存存储作为演示目的。


```bash
java -jar zipkin.jar
```

您可以通过访问 http://localhost:9411 验证安装效果

或者使用 Docker :

See [/dubbo-go-sample/zipkin/docker/docker-compose.yml](docker-compose.yml)

```dockerfile
version: '2.4'

services:
  # The zipkin process services the UI, and also exposes a POST endpoint that
  # instrumentation can send trace data to.
  zipkin:
    image: ghcr.io/openzipkin/zipkin-slim:${TAG:-latest}
    container_name: zipkin
    # Environment settings are defined here https://github.com/openzipkin/zipkin/blob/master/zipkin-server/README.md#environment-variables
    environment:
      - STORAGE_TYPE=mem
      # Point the zipkin at the storage backend
      - MYSQL_HOST=mysql
      # Uncomment to enable self-tracing
      # - SELF_TRACING_ENABLED=true
      # Uncomment to increase heap size
      # - JAVA_OPTS=-Xms128m -Xmx128m -XX:+ExitOnOutOfMemoryError
    ports:
      # Port used for the Zipkin UI and HTTP Api
      - 9411:9411
    # Uncomment to enable debug logging
    # command: --logging.level.zipkin2=DEBUG
```

## 如何运行
请参阅根目录中的 [HOWTO.md](../HOWTO_zh.md) 来运行本例。

