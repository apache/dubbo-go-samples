package main

import (
	"fmt"
	"github.com/apache/dubbo-go-samples/rpc/grpc/proto/greetgrpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	server := grpc.NewServer()
	greetgrpc.RegisterGreetServiceServer(server, &GreetGrpcServer{})
	listen, err := net.Listen("tcp", ":20000")
	if err != nil {
		panic(err)
	}
	fmt.Println("server start")
	err1 := server.Serve(listen)
	if err1 != nil {
		panic(err1)
	}
}
