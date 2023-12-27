package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo/client/service"

	"github.com/dubbogo/gost/log/logger"
)

var exceptionService = new(service.ExceptionService)

func main() {
	config.SetConsumerService(exceptionService)
	if err := config.Load(); err != nil {
		panic(err)
	}

	logger.Info("start to test dubbo")

	err := exceptionService.Throwable(context.Background(), "dubbo-go")
	if err != nil {
		logger.Error(err)
	}
}
