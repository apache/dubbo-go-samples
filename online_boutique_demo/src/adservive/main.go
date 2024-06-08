package main

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/apache/dubbo-go-demo/adservice/handler"
	hipstershop "github.com/apache/dubbo-go-demo/adservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := hipstershop.RegisterAdServiceHandler(srv, &handler.AdTripleService{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
