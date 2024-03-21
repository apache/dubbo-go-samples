# dubbogo-java

## Contents

- protobuf: 使用 proto 文件的结构体定义
- server
- client

请注意，该样例使用dubbo-go 3.2.0-rc1编写
我们测试的组合包括:

- [x] java-client -> dubbogo-server
- [x] java-server -> dubbogo-client

## 运行
1. 启动服务端
   - 使用 goland 启动 triple/gojava-go-server
   - 在 java-server 文件夹下执行 `sh run.sh` 启动 java server
2. 启动客户端
   - 使用 goland 启动 triple/gojava-go-client
   - 在 java-client 文件夹下执行 `sh run.sh` 启动 java client

## 注意
1. 接口命名须一致
   - java-server: GreeterImpl
   - go-client: 在conf中应类似如下定义
   ```yml
     Consumer:
       services:
         GreeterConsumer:
           # interface is for registry
           interface: org.apache.dubbo.sample.GreeterImpl
   ```
      