package main

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/apache/dubbo-go/common/logger"
    "github.com/apache/dubbo-go/config"

    _ "github.com/apache/dubbo-go/protocol/dubbo"
    _ "github.com/apache/dubbo-go/registry/protocol"

    _ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
    _ "github.com/apache/dubbo-go/filter/filter_impl"

    _ "github.com/apache/dubbo-go/cluster/cluster_impl"
    _ "github.com/apache/dubbo-go/cluster/loadbalance"
    // _ "github.com/apache/dubbo-go/metadata/service/remote"
    _ "github.com/apache/dubbo-go/metadata/service/inmemory"
    _ "github.com/apache/dubbo-go/registry/zookeeper"

    "github.com/apache/dubbo-go-samples/game/go-server-gate/pkg"
)

func main() {
    config.Load()

    go startHttp()

    initSignal()
}

func initSignal() {
    signals := make(chan os.Signal, 1)

    signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
    for {
        sig := <-signals
        logger.Infof("get signal %s", sig.String())
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

func startHttp() {
    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        res, err := pkg.Message(context.TODO(), "abc", "hello")
        if err != nil {
            _, _ = w.Write([]byte(err.Error()))
            return
        }

        b, err := json.Marshal(res)
        if err != nil {
            _, _ = w.Write([]byte(err.Error()))
            return
        }

        _, _ = w.Write(b)
    })

    _ = http.ListenAndServe(":8000", nil)
}