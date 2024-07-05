package main

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/recommendationservice/handler"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/recommendationservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("recommendationservice"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20011),
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

	if err := pb.RegisterRecommendationServiceHandler(srv, &handler.RecommendationService{}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
