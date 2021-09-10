package main

import (
	"context"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

type GreeterProvider struct {
	api.GreeterProviderBase
}

func main() {
	config.SetProviderService(&GreeterProvider{})
	config.Load()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	logger.SetLoggerLevel("warn")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Info("hello dubbogo this is info log")
			logger.Debug("hello dubbogo this is debug log")
			logger.Warn("hello dubbogo this is warn log")
			time.Sleep(time.Second * 1)
		}
	}
}
