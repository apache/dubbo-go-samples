package main

import (
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/handler"
	hipstershop "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"
	"github.com/dubbogo/gost/log/logger"
	_ "github.com/dubbogo/gost/log/logger"
)

func main() {
	//server
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
	)
	if err != nil {
		panic(err)
	}

	if err := hipstershop.RegisterCartServiceHandler(srv, &handler.CartService{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}

}
