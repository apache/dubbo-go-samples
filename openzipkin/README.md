
### Install Zipkin

First Way:

Follow [Zipkin's quick start](https://zipkin.io/pages/quickstart.html) to install zipkin.

```bash
curl -sSL https://zipkin.io/quickstart.sh | bash -s
```

Zipkin supports various backend storages including Cassandra, ElasticSearch and MySQL. Here we use the simplest storage - in-memory for demo purpose.

```bash
java -jar zipkin.jar
```

Once the process starts, you can verify zipkin server works by access http://localhost:9411

Or Use Docker:

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


```bash
docker-compose -f docker/docker-compose.yml up -d zipkin
```



![img.png](img.png)




