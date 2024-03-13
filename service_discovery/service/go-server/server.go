package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3"
	// "dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	greet "github.com/apache/dubbo-go-samples/service_discovery/service/proto"
	"github.com/dubbogo/gost/log/logger"
)

type GreetTripleServer struct {
}

func (svr *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func (s *GreetTripleServer) MethodMapper() map[string]string {
	return map[string]string{
		"Greet":    "greet",
		"SayHello": "sayHello",
	}
}

func (svr *GreetTripleServer) SayHello(ctx context.Context, req *greet.SayHelloRequest) (*greet.SayHelloResponse, error) {
	resp := &greet.SayHelloResponse{Hello: req.Name}
	return resp, nil
}

func main() {

	ins, err := dubbo.NewInstance(
		dubbo.WithName("dubbo-go-server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress("127.0.0.1:8848"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20022),
		),
	)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err := greet.RegisterGreetServiceHandler(srv, &GreetTripleServer{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}

}
