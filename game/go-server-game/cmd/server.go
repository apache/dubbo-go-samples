package main

import (
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/apache/dubbo-go/common/logger"
    "github.com/apache/dubbo-go/config"
    hessian "github.com/apache/dubbo-go-hessian2"

    _ "github.com/apache/dubbo-go/protocol/dubbo"
    _ "github.com/apache/dubbo-go/registry/protocol"

    _ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
    _ "github.com/apache/dubbo-go/filter/filter_impl"

    _ "github.com/apache/dubbo-go/cluster/cluster_impl"
    _ "github.com/apache/dubbo-go/cluster/loadbalance"
    // _ "github.com/apache/dubbo-go/metadata/service/remote"
    _ "github.com/apache/dubbo-go/metadata/service/inmemory"
    _ "github.com/apache/dubbo-go/registry/zookeeper"

    "github.com/apache/dubbo-go-samples/game/go-server-game/pkg"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

func init() {
    config.SetProviderService(new(pkg.BasketballService))

    config.SetConsumerService(pkg.GateBasketball)

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
