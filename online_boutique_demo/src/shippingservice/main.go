package main

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/shippingservice/handler"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/shippingservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

var (
	name    = "shippingservice"
	version = "1.0.0"
)

func main() {

	ins, err := dubbo.NewInstance(
		dubbo.WithName(name),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
	)
	if err != nil {
		panic(err)
	}
	//server
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}

	if err := pb.RegisterShippingServiceHandler(srv, new(handler.ShippingService)); err != nil {
		logger.Fatal(err)
	}

	// Run service
	if err := srv.Serve(); err != nil {
		logger.Fatal(err)
	}
}
