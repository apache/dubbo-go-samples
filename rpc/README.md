# Multi-protocol support

## 1.Supports Grpc protocol

​	It's worth noting that triple works naturally with Grpc.

Implementation steps



​	1. Create a Grpc client or server，For details, please refer to[Documentation | gRPC](https://grpc.io/docs/)

Run the instance code client

```shell
go run ./grpc/client/client.go
```

Run the instance code server

```shell
go run ./grpc/server/server.go
```



​	2. Create a Triple client or server

Run the instance code client

```shell
go run ./triple/client/client.go
```

Run the instance code server

```shell
go run ./triple/server/server.go
```



3. Example Results（The results include Unary ClientStream ServerStream BidiStream ）

   3.1 The client using Triple calls the server of Grpc

   Results：

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

   

   3.2 The client using Grpc calls the server of Triple

   Results：

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