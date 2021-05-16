# Multi Zone Example

## Background

Sometimes it is necessary to deploy Dubbo services into different zones for the sake of disaster tolerance or high availability. In order to consumer the services among the different registry centers, you may need to consider configuring your consumer like this:

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

## Run Steps

1. Start Zookeeper in docker environment.

    ```bash
    make -f ../build/Makefile docker-up
    ```
   
2. Start Server

    Start Hangzhou zone server.
    ```bash
   cd go-server-hz
   make -f ../../build/Makefile clean start  
   ```
   
   Start Shanghai zone server.
   ```bash
   cd go-server-sh
   make -f ../../build/Makefile clean start  
  
3. Run Consumer

    ```bash
    cd go-client
    make -f ../../build/Makefile run
    ```
   
4. Cleanup

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