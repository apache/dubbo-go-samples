package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	health "dubbo.apache.org/dubbo-go/v3/protocol/triple/health/triple_health"
	hipstershop "github.com/apache/dubbo-go-demo/adservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}
	checkhealth(cli)
	svc, err := hipstershop.NewAdService(cli)
	if err != nil {
		panic(err)
	}
	resp, err := svc.GetAds(context.Background(), &hipstershop.AdRequest{
		ContextKeys: []string{"Cookie"},
	})
	if err != nil {
		panic(err)
	}
	logger.Infof("get resp: %v", resp)
}

func checkhealth(cli *client.Client) error {
	svc, err := health.NewHealth(cli)
	if err != nil {
		panic(err)
	}
	check, err := svc.Check(context.Background(), &health.HealthCheckRequest{Service: hipstershop.AdServiceName})
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info(hipstershop.AdServiceName, "'s health", check.String())
	}
	watch, err := svc.Watch(context.Background(), &health.HealthCheckRequest{Service: hipstershop.AdServiceName})
	if err != nil {
		logger.Error(err)
	} else {
		if watch.Recv() {
			logger.Info(hipstershop.AdServiceName, "greet.GreetService's health", watch.Msg().String())
		}
	}
	return nil
}
