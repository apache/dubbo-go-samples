# 泛化调用

泛化接口调用方式主要用于客户端没有 API 接口及模型类元的情况，参数及返回值中的所有 POJO 均用 Map 表示，通常用于框架集成，比如：实现一个通用的服务测试框架，可通过`GenericService`调用所有服务实现。更多信息请参考文档。

## 开始

### 使用方法

1. 启动 zookeeper

   ```shell
   cd ./default/go-server/docker \
     && docker-compose up -d
   ```

2. 启动服务端运行提供者

   1. go

      使用 goland 启动 generic-default-go-server

   2. java

      在goland启动generic-default-java-server
      
      或

      在 java-server 文件夹下执行 `sh run.sh` 启动 java server

3. 启动客户端运行消费者发启泛化调用

   1. go

      使用 goland 启动 generic-default-go-client

   2. java

      在goland启动generic-default-java-client

      或

      在 java-client 文件夹下执行 `sh run.sh` 启动 java client