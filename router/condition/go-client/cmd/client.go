package main

import (
	"context"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/direct/proto"
)

const (
	RegistryAddress = "127.0.0.1:8848"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("condition-client"),
		dubbo.WithRegistry(
			registry.WithNacos(),
			registry.WithAddress(RegistryAddress),
		),
		dubbo.WithConfigCenter( // configure config center to enable condition router
			config_center.WithNacos(),
			config_center.WithAddress(RegistryAddress),
		),
	)

	if err != nil {
		logger.Errorf("new instance failed: %v", err)
		panic(err)
	}

	cli, err := ins.NewClient()

	if err != nil {
		logger.Errorf("new client failed: %v", err)
		panic(err)
	}

	srv, err := greet.NewGreetService(cli)

	if err != nil {
		logger.Errorf("new service failed: %v", err)
		panic(err)
	}

	for {
		time.Sleep(5 * time.Second) // sleep 5 seconds
		rep, err := srv.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
		printRes(rep, err)
	}

}

func printRes(rep *greet.GreetResponse, err error) {
	if err != nil {
		logger.Errorf("call greet method failed: %v", err)
	} else {
		logger.Infof("receive: %s", rep.GetGreeting())
	}
}
