# Dubbo Streaming Sample

## 1. 介绍

本案例演示如何在 Dubbo 中使用流式通信，包括：
- Go 语言的流式通信实现
- Java 语言的流式通信实现
- Go 与 Java 之间的互通性验证

支持的流式模式：
- 一元调用 (Unary): 单个请求，单个响应
- 双向流 (Bidirectional Stream): 多个请求，多个响应
- 客户端流 (Client Stream): 多个请求，单个响应
- 服务端流 (Server Stream): 单个请求，多个响应

## 2. Proto 定义

在 proto 文件中定义流式通信方法，在需要流式传输的参数前添加 `stream` 关键字：

```protobuf
service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
  rpc GreetStream(stream GreetStreamRequest) returns (stream GreetStreamResponse) {}
  rpc GreetClientStream(stream GreetClientStreamRequest) returns (GreetClientStreamResponse) {}
  rpc GreetServerStream(GreetServerStreamRequest) returns (stream GreetServerStreamResponse) {}
}
```

## 3. Go 语言实现

### 3.1 Go 服务端

源文件路径: `streaming/go-server/cmd/server.go`

```go
type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func (srv *GreetTripleServer) GreetStream(ctx context.Context, stream greet.GreetService_GreetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if triple.IsEnded(err) {
				break
			}
			return fmt.Errorf("triple BidiStream recv error: %s", err)
		}
		if err := stream.Send(&greet.GreetStreamResponse{Greeting: req.Name}); err != nil {
			return fmt.Errorf("triple BidiStream send error: %s", err)
		}
	}
	return nil
}

func (srv *GreetTripleServer) GreetClientStream(ctx context.Context, stream greet.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
	var reqs []string
	for stream.Recv() {
		reqs = append(reqs, stream.Msg().Name)
	}
	if stream.Err() != nil && !triple.IsEnded(stream.Err()) {
		return nil, fmt.Errorf("triple ClientStream recv err: %s", stream.Err())
	}
	resp := &greet.GreetClientStreamResponse{
		Greeting: strings.Join(reqs, ","),
	}

	return resp, nil
}

func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greet.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: req.Name}); err != nil {
			return fmt.Errorf("triple ServerStream send err: %s", err)
		}
	}
	return nil
}
```

### 3.2 Go 客户端

源文件路径: `streaming/go-client/cmd/client.go`

```go
func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}
	TestClient(svc)
}

func TestClient(cli greet.GreetService) {
	if err := testUnary(cli); err != nil {
		logger.Error(err)
	}

	if err := testBidiStream(cli); err != nil {
		logger.Error(err)
	}

	if err := testClientStream(cli); err != nil {
		logger.Error(err)
	}

	if err := testServerStream(cli); err != nil {
		logger.Error(err)
	}
}

func testUnary(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE unary call")
	resp, err := cli.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE unary call resp: %s", resp.Greeting)
	return nil
}

func testBidiStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE bidi stream")
	stream, err := cli.GreetStream(context.Background())
	if err != nil {
		return err
	}
	if sendErr := stream.Send(&greet.GreetStreamRequest{Name: "triple"}); sendErr != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE bidi stream resp: %s", resp.Greeting)
	if err := stream.CloseRequest(); err != nil {
		return err
	}
	if err := stream.CloseResponse(); err != nil {
		return err
	}
	return nil
}

func testClientStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE client stream")
	stream, err := cli.GreetClientStream(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		if sendErr := stream.Send(&greet.GreetClientStreamRequest{Name: "triple"}); sendErr != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE client stream resp: %s", resp.Greeting)
	return nil
}

func testServerStream(cli greet.GreetService) error {
	logger.Info("start to test TRIPLE server stream")
	stream, err := cli.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "triple"})
	if err != nil {
		return err
	}
	for stream.Recv() {
		logger.Infof("TRIPLE server stream resp: %s", stream.Msg().Greeting)
	}
	if stream.Err() != nil {
		return err
	}
	if err := stream.Close(); err != nil {
		return err
	}
	return nil
}
```

## 4. Java 语言实现

### 4.1 项目结构

```
streaming/
├── pom.xml                    # 父 POM
├── proto/                     # 共享 Proto 文件
│   └── greet.proto
├── java-server/              # Java 服务端
│   ├── pom.xml
│   ├── src/main/java/
│   │   └── org/apache/dubbo/samples/tri/streaming/
│   │       ├── StreamingServer.java
│   │       └── GreeterImpl.java
│   └── run.sh
└── java-client/              # Java 客户端
    ├── pom.xml
    ├── src/main/java/
    │   └── org/apache/dubbo/samples/tri/streaming/
    │       └── StreamingClient.java
    └── run.sh
```

### 4.2 Java 服务端

源文件路径: `streaming/java-server/src/main/java/org/apache/dubbo/samples/tri/streaming/GreeterImpl.java`

实现了所有四种 RPC 模式：

```java
public class GreeterImpl extends DubboGreetServiceTriple.GreetServiceImplBase {
    
    // 一元调用
    @Override
    public GreetResponse greet(GreetRequest request) {
        return GreetResponse.newBuilder()
                .setGreeting("Hello " + request.getName())
                .build();
    }
    
    // 双向流
    @Override
    public StreamObserver<GreetStreamRequest> greetStream(
            StreamObserver<GreetStreamResponse> responseObserver) {
        // 实现双向流逻辑
    }
    
    // 客户端流
    @Override
    public StreamObserver<GreetClientStreamRequest> greetClientStream(
            StreamObserver<GreetClientStreamResponse> responseObserver) {
        // 实现客户端流逻辑
    }
    
    // 服务端流
    @Override
    public void greetServerStream(GreetServerStreamRequest request,
            StreamObserver<GreetServerStreamResponse> responseObserver) {
        // 发送 10 个响应
        for (int i = 0; i < 10; i++) {
            responseObserver.onNext(response);
        }
        responseObserver.onCompleted();
    }
}
```

### 4.3 Java 客户端

源文件路径: `streaming/java-client/src/main/java/org/apache/dubbo/samples/tri/streaming/StreamingClient.java`

测试所有流式模式，输出格式清晰易读。

## 5. 运行示例

### 5.1 运行 Go 服务端和客户端

```bash
# 启动 Go 服务端
cd streaming/go-server
go run cmd/server.go

# 在另一个终端启动 Go 客户端
cd streaming/go-client
go run cmd/client.go
```

### 5.2 运行 Java 服务端和客户端

```bash
# 编译父项目
cd streaming
mvn clean install

# 启动 Java 服务端
cd java-server
./run.sh
# 或者
mvn compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.tri.streaming.StreamingServer"

# 在另一个终端启动 Java 客户端
cd java-client
./run.sh
# 或者
mvn compile exec:java -Dexec.mainClass="org.apache.dubbo.samples.tri.streaming.StreamingClient"
```

## 6. 案例效果

### 6.1 Java 客户端输出示例

```
======================================================================
 Starting Dubbo Streaming Client Tests
======================================================================
 Connected to server: tri://127.0.0.1:20000
======================================================================

======================================================================
 TEST 1: Bidirectional Streaming
======================================================================
    Sending request #0: Client-0
    Received response #1: Echo from biStream: Client-0
    Sending request #1: Client-1
    Received response #2: Echo from biStream: Client-1
    Sending request #2: Client-2
    Received response #3: Echo from biStream: Client-2
    Sending request #3: Client-3
    Received response #4: Echo from biStream: Client-3
    Sending request #4: Client-4
    Received response #5: Echo from biStream: Client-4

   All requests sent, waiting for responses...

   BiStream completed - Received 5 responses

======================================================================
 TEST 2: Server Streaming
======================================================================
    Sending request: StreamingClient
    Waiting for server stream responses...

    Received response #1: Response 0 from serverStream for StreamingClient
    Received response #2: Response 1 from serverStream for StreamingClient
  ...
    Received response #10: Response 9 from serverStream for StreamingClient

   ServerStream completed - Received 10 responses

======================================================================
 TEST RESULTS SUMMARY
======================================================================
  Bidirectional Streaming:  PASSED
  Server Streaming:  PASSED
----------------------------------------------------------------------
   Overall: ALL TESTS PASSED!
======================================================================
```

### 6.2 Go 客户端输出示例

```
INFO    cmd/client.go:69    start to test TRIPLE unary call
INFO    cmd/client.go:74    TRIPLE unary call resp: Hello triple

INFO    cmd/client.go:79    start to test TRIPLE bidi stream
INFO    cmd/client.go:91    TRIPLE bidi stream resp: Echo from biStream: triple

INFO    cmd/client.go:102   start to test TRIPLE client stream
INFO    cmd/client.go:116   TRIPLE client stream resp: Received 5 names: triple, triple, triple, triple, triple

INFO    cmd/client.go:121   start to test TRIPLE server stream
INFO    cmd/client.go:127   TRIPLE server stream resp: Response 0 from serverStream for triple
INFO    cmd/client.go:127   TRIPLE server stream resp: Response 1 from serverStream for triple
...
INFO    cmd/client.go:127   TRIPLE server stream resp: Response 9 from serverStream for triple
```

## 注意

由于Golang-server和Java-server同时监听localhost:20000端口,你不能同时启动两个服务端


