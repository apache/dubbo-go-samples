package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

func init() {
	config.SetProviderService(new(BasketballService))

	config.SetConsumerService(GateBasketball)

	hessian.RegisterPOJO(&pojo.Result{})
}

func main() {
	config.Load()

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %#s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			logger.Infof("app need reload")
		default:
			time.AfterFunc(time.Duration(time.Second*5), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			// The program exits normally or timeout forcibly exits.
			logger.Warnf("app exit now...")
			return
		}
	}
}