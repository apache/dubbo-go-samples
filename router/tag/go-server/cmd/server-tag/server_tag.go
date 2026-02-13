package main

import (
	"context"
	"strings"

	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	greet "github.com/apache/dubbo-go-samples/router/tag/proto"

	"github.com/dubbogo/gost/log/logger"
)

const (
	RegistryAddress = "127.0.0.1:8848"
)

type GreetServer struct {
	srvName string
}

func (srv *GreetServer) Greet(_ context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	rep := &greet.GreetResponse{Greeting: "receive: " + req.Name + ", response from: " + srv.srvName}
	return rep, nil
}

func main() {

	ins, err := dubbo.NewInstance(
		dubbo.WithName("tag-server"),
		dubbo.WithTag("test-tag"), // set application's tag
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(RegistryAddress),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20001),
		),
	)

	if err != nil {
		logger.Errorf("new dubbo instance failed: %v", err)
		panic(err)
	}

	srv, err := ins.NewServer()

	if err != nil {
		logger.Errorf("new dubbo server failed: %v", err)
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetServer{srvName: "server-with-tag"}); err != nil {
		logger.Errorf("register greet handler failed: %v", err)
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Errorf("server serve failed: %v", err)
		if strings.Contains(err.Error(), "client not connected") {
			logger.Errorf("hint: Nacos client not connected (gRPC). Check %s is reachable and gRPC port %d is open (Nacos 2.x default).", RegistryAddress, 20001)
		}
		panic(err)
	}

}
