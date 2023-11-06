package main

import (
	"context"
	greet "github.com/apache/dubbo-go-samples/rpc/grpc/proto"
	"github.com/apache/dubbo-go-samples/rpc/grpc/proto/greetgrpc"
	"github.com/dubbogo/gost/log/logger"
	"google.golang.org/grpc"
	"io"
)

func main() {
	dial, err := grpc.Dial("127.0.0.1:20000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := greetgrpc.NewGreetServiceClient(dial)
	TestGrpcClient(client)
}

func TestGrpcClient(gsc greetgrpc.GreetServiceClient) {
	if err := unaryCall(gsc); err != nil {
		logger.Error(err)
	}

	if err := clientStreamCall(gsc); err != nil {
		logger.Error(err)
	}

	if err := serverStreamCall(gsc); err != nil {
		logger.Error(err)
	}

	if err := bidiStreamCall(gsc); err != nil {
		logger.Error(err)
	}
}

func unaryCall(gsc greetgrpc.GreetServiceClient) error {
	logger.Info("start to test Grpc unary call")
	response, err := gsc.Greet(context.Background(), &greet.GreetRequest{Name: "Grpc"})
	if err != nil {
		return err
	}
	logger.Infof("Grpc unary call resp: %s", response.Greeting)
	return nil
}

func clientStreamCall(gsc greetgrpc.GreetServiceClient) error {
	logger.Info("start to test Grpc client stream")
	stream, err := gsc.GreetClientStream(context.Background())
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		if sendErr := stream.Send(&greet.GreetClientStreamRequest{Name: "Grpc"}); sendErr != nil {
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	logger.Infof("Grpc client stream resp: %s", resp.Greeting)
	return nil
}

func serverStreamCall(gsc greetgrpc.GreetServiceClient) error {
	logger.Info("start to test Grpc server stream")
	stream, err := gsc.GreetServerStream(context.Background(), &greet.GreetServerStreamRequest{Name: "Grpc"})
	if err != nil {
		return err
	}
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		logger.Infof("Grpc server stream resp: %s", recv.Greeting)
	}
	return nil
}

func bidiStreamCall(gsc greetgrpc.GreetServiceClient) error {
	logger.Info("start to test Grpc bidi stream")
	stream, err := gsc.GreetStream(context.Background())
	if err != nil {
		return err
	}
	if sendErr := stream.Send(&greet.GreetStreamRequest{Name: "Grpc"}); sendErr != nil {
		return err
	}
	resp, err := stream.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Grpc bidi stream resp: %s", resp.Greeting)
	return nil
}
