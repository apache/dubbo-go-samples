# 多协议支持

## 1.支持Grpc协议

​	值得说明的是，triple协议是天然能与Grpc协议互通的。

实现步骤



​	1. 创建Grpc的客户端或者服务端，详情请参考[Documentation | gRPC](https://grpc.io/docs/)

运行实例代码客户端

```shell
go run ./grpc/client/client.go
```

运行实例代码服务端

```shell
go run ./grpc/server/server.go
```



 	2. 创建Triple的客户端或者服务端

运行实例代码客户端

```shell
go run ./triple/client/client.go
```

运行实例代码服务端

```shell
go run ./triple/server/server.go
```



3. 示例结果（结果中包括 Unary ClientStream ServerStream BidiStream 四种调用）

   3.1 使用Triple的客户端调用Grpc的服务端

   结果：

   ```
   2023-11-06 10:12:26     INFO    logger/logging.go:22    [start to test TRIPLE unary call]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE unary call resp: [Grpc greet server receive triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:22    [start to test TRIPLE bidi stream]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE bidi stream resp: [Grpc greetStream server receive triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:22    [start to test TRIPLE client stream]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE client stream resp: [Grpc greetClientStream server receive triple,triple,triple,triple,triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:22    [start to test TRIPLE server stream]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE server stream resp: [Grpc greetServerStream server receive triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE server stream resp: [Grpc greetServerStream server receive triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE server stream resp: [Grpc greetServerStream server receive triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE server stream resp: [Grpc greetServerStream server receive triple]
   2023-11-06 10:12:26     INFO    logger/logging.go:42    TRIPLE server stream resp: [Grpc greetServerStream server receive triple]
   
   ```

   3.2 使用Grpc的客户端调用Triple的服务端

   结果：

   ```
   2023-11-06T10:18:52.177+0800    INFO    client/client.go:40     start to test Grpc unary call
   2023-11-06T10:18:52.184+0800    INFO    client/client.go:45     Grpc unary call resp: Triple greet server receive Grpc
   2023-11-06T10:18:52.184+0800    INFO    client/client.go:50     start to test Grpc client stream
   2023-11-06T10:18:52.185+0800    INFO    client/client.go:64     Grpc client stream resp: Triple greetClientStream server receive Grpc,Grpc,Grpc,Grpc,Grpc
   2023-11-06T10:18:52.185+0800    INFO    client/client.go:69     start to test Grpc server stream
   2023-11-06T10:18:52.186+0800    INFO    client/client.go:82     Grpc server stream resp: Triple greetServerStream server receive Grpc
   2023-11-06T10:18:52.186+0800    INFO    client/client.go:82     Grpc server stream resp: Triple greetServerStream server receive Grpc
   2023-11-06T10:18:52.186+0800    INFO    client/client.go:82     Grpc server stream resp: Triple greetServerStream server receive Grpc
   2023-11-06T10:18:52.187+0800    INFO    client/client.go:82     Grpc server stream resp: Triple greetServerStream server receive Grpc
   2023-11-06T10:18:52.187+0800    INFO    client/client.go:82     Grpc server stream resp: Triple greetServerStream server receive Grpc
   2023-11-06T10:18:52.187+0800    INFO    client/client.go:88     start to test Grpc bidi stream
   2023-11-06T10:18:52.187+0800    INFO    client/client.go:100    Grpc bidi stream resp: Triple greetStream server receive Grpc
   
   ```

   

​	