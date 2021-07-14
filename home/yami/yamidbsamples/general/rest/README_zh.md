## Dubbo-Go Restful 示例

[Official Document](https://dubbo.apache.org/zh/blog/2021/01/14/dubbo-go-%E4%B8%AD-rest-%E5%8D%8F%E8%AE%AE%E5%AE%9E%E7%8E%B0/)

### 1. 介绍

使用Dubbo-go开发restful应用.

### 2. 使用示例

生产者

```yaml
# server.yml

# service config
services:
  "UserProvider":
    registry: "demoZk"
    # using rest protocl
    protocol : "rest"
    interface : "org.apache.dubbo.UserProvider"
    loadbalance: "random"
    warmup: "100"
    cluster: "failover"
    # the http path
    rest_path: "/UserProvider"
    methods:
      - name: "GetUser"
        # the http path
        rest_path: "/GetUser"
        # the http method
        rest_method: "GET"
        # the query param 'id' mapping to the first parameter in this method
        rest_query_params: "0:id"
      - name: "GetUser0"
        rest_path: "/GetUser0/{id}"
        rest_method: "POST"
        # the path param 'name' mapping to the second parameter in this method,
        # and the path param 'age' mapping to the third parameter.
        rest_query_params: "1:name,2:age"
        # the path param 'id' mapping to the first parameter in this method
        rest_path_params: "0:id"
        rest_produces: "application/json"
        rest_consumes: "application/json;charset=utf-8,*/*"
      - name: "GetUser3"
        rest_path: "/GetUser3"
        rest_method: "GET"
      - name: "GetUsers"
        rest_path: "/GetUsers"
        rest_method: "POST"
        rest_body: 0
      - name: "GetUser1"
        rest_path: "/GetUser1"
        rest_method: "GET"

protocols:
  "rest":
    name: "rest"
    ip: "0.0.0.0"
    port: 8888
```

消费者

```yaml
# client.yml

# reference config
references:
  "UserProvider":
    registry: "demoZk"
    protocol : "rest"
    interface : "org.apache.dubbo.UserProvider"
    cluster: "failover"
    rest_path: "/UserProvider"
    methods:
      - name: "GetUser"
        rest_path: "/GetUser"
        rest_method: "GET"
        rest_query_params: "0:id"
      - name: "GetUser0"
        rest_path: "/GetUser0/{id}"
        rest_method: "POST"
        rest_query_params: "1:name,2:age"
        rest_path_params: "0:id"
      - name: "GetUser3"
        rest_path: "/GetUser3"
        rest_method: "GET"
      - name: "GetUsers"
        rest_path: "/GetUsers"
        rest_method: "POST"
        rest_body: 0
      - name: "GetUser1"
        rest_path: "/GetUser1"
        rest_method: "GET"
```