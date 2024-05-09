# Dubbo java and go interoperability, with service discovery ands protocol

This example shows dubbo-go's service discovery and java-go interoperation feature with Nacos as registry.

> before run the code, you should Follow the instruction to <a href="https://dubbo-next.staged.apache.org/zh-cn/overview/reference/integrations/nacos/" target="_blank">install and start Nacos server</a>.

本示例针对 Dubbo3 版本用户，演示如何基于 Dubbo3 应用级服务发现实现 Java 与 Go 体系互调。

## dubbo java 调用 dubbo go

1. 启动 go server

    ```shell
    go run ./go-server/cmd/server.go
    ```

2. 启动 java client

    ```shell
    ./java-client/run.sh
    ```

## dubbo go 调用 dubbo java

> 本示例当前存在问题，请跟进排查：https://github.com/apache/dubbo-go/issues/2653

1. 启动 java server

    ```shell
   ./java-server/run.sh
    ```

2. 启动 go client

    ```shell
    go run ./go-client/cmd/client.go
    ```