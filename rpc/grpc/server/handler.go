package main

import (
	"context"
	"fmt"
	greet "github.com/apache/dubbo-go-samples/rpc/grpc/proto"
	"github.com/apache/dubbo-go-samples/rpc/grpc/proto/greetgrpc"
	"io"
	"strings"
)

type GreetGrpcServer struct {
	greetgrpc.UnimplementedGreetServiceServer
}

func (receiver *GreetGrpcServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: "Grpc greet server receive " + req.Name}
	return resp, nil
}

func (receiver *GreetGrpcServer) GreetStream(stream greetgrpc.GreetService_GreetStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Grpc BidiStream recv error: %s", err)
		}
		if err := stream.Send(&greet.GreetStreamResponse{Greeting: "Grpc greetStream server receive " + req.Name}); err != nil {
			return fmt.Errorf("Grpc BidiStream send error: %s", err)
		}
	}
	return nil
}

func (receiver *GreetGrpcServer) GreetClientStream(stream greetgrpc.GreetService_GreetClientStreamServer) error {
	var reqs []string
	for {
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Grpc ClientStream recv err: %s", err.Error())
		}
		reqs = append(reqs, recv.Name)
	}
	resp := &greet.GreetClientStreamResponse{
		Greeting: "Grpc greetClientStream server receive " + strings.Join(reqs, ","),
	}
	stream.SendAndClose(resp)
	return nil
}

func (receiver *GreetGrpcServer) GreetServerStream(req *greet.GreetServerStreamRequest, stream greetgrpc.GreetService_GreetServerStreamServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&greet.GreetServerStreamResponse{Greeting: "Grpc greetServerStream server receive " + req.Name}); err != nil {
			return fmt.Errorf("Grpc ServerStream send err: %s", err)
		}
	}
	return nil
}
