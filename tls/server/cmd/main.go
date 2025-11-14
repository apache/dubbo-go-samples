package main

import (
	"context"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"dubbo.apache.org/dubbo-go/v3/tls"

	"github.com/dubbogo/gost/log/logger"
)

import (
	"github.com/apache/dubbo-go-samples/tls/proto"
)

type GreetTripleServer struct {
}

func (srv *GreetTripleServer) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {
	resp := &greet.GreetResponse{Greeting: req.Name}
	return resp, nil
}

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
		server.WithServerTLSOption(
			tls.WithCertFile("../../x509/server2_cert.pem"),
			tls.WithKeyFile("../../x509/server2_key.pem"),
			tls.WithServerName("dubbogo.test.example.com"),
		),
	)
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
