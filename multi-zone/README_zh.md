# 多注册中心示例

## 背景

有些时候为了容灾或者高可用的目的，有必要把 Dubbo 服务部署到不同的区域去。为了能够同时消费到注册在不同注册中心上的服务，需要按照以下的方式来配置你的消费方：

```yaml
# registry config
registries:
  "shanghaiZk":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2182"
    username: ""
    password: ""
    zone: shanghai
    weight: 50

  "hangzhouZk":
    protocol: "zookeeper"
    timeout: "3s"
    address: "127.0.0.1:2183"
    zone: hangzhou
    weight: 200

# reference config
references:
  "UserProvider":
    registry: "shanghaiZk,hangzhouZk"
    protocol: "dubbo"
    interface: "org.apache.dubbo.UserProvider"
    cluster: "failover"
    methods:
      - name: "GetUser"
        retries: 3
```

## 运行步骤

1. 在 Docker 环境中启动 Zookeeper 服务器

    ```bash
    make -f ../build/Makefile docker-up
    ```
   
2. 启动服务端

   启动杭州区域的服务器。
    ```bash
   cd go-server-hz
   make -f ../../build/Makefile clean start  
   ```
   
   启动上海区域的服务器。
   ```bash
   cd go-server-sh
   make -f ../../build/Makefile clean start  
  
3. 运行消费方

    ```bash
    cd go-client
    make -f ../../build/Makefile run
    ```
   
4. 清理

   ```bash
   cd go-server-hz && \
        make -f ../../build/Makefile clean && \
        cd ..
   cd go-server-sh && \
        make -f ../../build/Makefile clean && \
        cd ..
   cd go-client && \
        make -f ../../build/Makefile clean && \
        cd ..
   make -f ../build/Makefile docker-down
   ```