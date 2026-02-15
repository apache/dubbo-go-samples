package main

import (
	"context"
	"strconv"
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	RegistryAddress = "127.0.0.1:8848"
	TriPort         = 20001
)

type GreetServer struct {
	srvPort int
}

func (srv *GreetServer) Greet(_ context.Context, req *greet.GreetRequest) (rep *greet.GreetResponse, err error) {
	rep = &greet.GreetResponse{Greeting: req.Name + "from: " + strconv.Itoa(srv.srvPort)}
	return rep, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("condition-server"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(RegistryAddress),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(TriPort),
		),
	)

	if err != nil {
		logger.Errorf("new instance failed: %v", err)
		panic(err)
	}

	srv, err := ins.NewServer()

	if err != nil {
		logger.Errorf("new server failed: %v", err)
		panic(err)
	}

	if err := greet.RegisterGreetServiceHandler(srv, &GreetServer{srvPort: TriPort}); err != nil {
		logger.Errorf("register service failed: %v", err)
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Errorf("server serve failed: %v", err)
		if strings.Contains(err.Error(), "client not connected") {
			logger.Errorf("hint: Nacos client not connected (gRPC). Check %s is reachable and gRPC port %d is open (Nacos 2.x default).", RegistryAddress, TriPort)
		}
		panic(err)
	}
}
