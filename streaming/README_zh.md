# Dubbo-go Streaming Sample

## 1.介绍

本案例演示如何在Dubbo-go中使用流式通信

## 2.如何使用Dubbo-go流式通信

在proto文件中需要流式通信的方法的参数前面添加stream，使用proto-gen-triple生成相应文件

```protobuf
service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {}
  rpc GreetStream(stream GreetStreamRequest) returns (stream GreetStreamResponse) {}
  rpc GreetClientStream(stream GreetClientStreamRequest) returns (GreetClientStreamResponse) {}
  rpc GreetServerStream(GreetServerStreamRequest) returns (stream GreetServerStreamResponse) {}
}
```

编写服务端handler文件

源文件路径: dubbo-go-sample/streaming/go-server/handler.go

```go
type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func (srv *GreetTripleServer) GreetStream(ctx context.Context, stream greettriple.GreetService_GreetStreamServer) error {
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

func (srv *GreetTripleServer) GreetClientStream(ctx context.Context, stream greettriple.GreetService_GreetClientStreamServer) (*greet.GreetClientStreamResponse, error) {
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

func (srv *GreetTripleServer) GreetServerStream(ctx context.Context, req *greet.GreetServerStreamRequest, stream greettriple.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: req.Name}); err != nil {
			return fmt.Errorf("triple ServerStream send err: %s", err)
		}
	}
	return nil
}
```

编写客户端client文件

源文件路径: dubbo-go-sample/streaming/go-client/client.go

```go
func main() {
	cli, err := client.NewClient(
		client.WithURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greettriple.NewGreetService(cli)
	if err != nil {
		panic(err)
	}
	TestClient(svc)
}

func TestClient(cli greettriple.GreetService) {
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

func testUnary(cli greettriple.GreetService) error {
	logger.Info("start to test TRIPLE unary call")
	resp, err := cli.Greet(context.Background(), &greet.GreetRequest{Name: "triple"})
	if err != nil {
		return err
	}
	logger.Infof("TRIPLE unary call resp: %s", resp.Greeting)
	return nil
}

func testBidiStream(cli greettriple.GreetService) error {
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

func testClientStream(cli greettriple.GreetService) error {
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

func testServerStream(cli greettriple.GreetService) error {
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

## 3.案例效果

运行服务端和客户端，可以看到请求正常返回

```
[start to test TRIPLE unary call]
TRIPLE unary call resp: [triple]
[start to test TRIPLE bidi stream]
TRIPLE bidi stream resp: [triple]
[start to test TRIPLE client stream]
TRIPLE client stream resp: [triple,triple,triple,triple,triple]
[start to test TRIPLE server stream]
TRIPLE server stream resp: [triple]
TRIPLE server stream resp: [triple]
TRIPLE server stream resp: [triple]
TRIPLE server stream resp: [triple]
TRIPLE server stream resp: [triple]
```

