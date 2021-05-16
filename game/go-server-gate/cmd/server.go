package main

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    hessian "github.com/apache/dubbo-go-hessian2"
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
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

func init() {
    config.SetProviderService(new(pkg.BasketballService))

    config.SetConsumerService(pkg.GameBasketball)

    hessian.RegisterPOJO(&pojo.Result{})
}


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

    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        res, err := pkg.Login(context.TODO(), r.URL.Query().Get("name"))
        if err != nil {
            responseWithOrigin(w, r, 200, []byte(err.Error()))
            return
        }

        b, err := json.Marshal(res)
        if err != nil {
            responseWithOrigin(w, r, 200, []byte(err.Error()))
            return
        }
        responseWithOrigin(w, r, 200, b)
    })

    http.HandleFunc("/score", func(w http.ResponseWriter, r *http.Request) {
        reqBody, err := ioutil.ReadAll(r.Body)
        if err != nil {
            logger.Error(err.Error())
        }
        var info pojo.Info
        json.Unmarshal(reqBody, &info)
        res, err := pkg.Score(context.TODO(), info.Name, strconv.Itoa(info.Score))
        if err != nil {
            responseWithOrigin(w, r, 200, []byte(err.Error()))
            return
        }

        b, err := json.Marshal(res)
        if err != nil {
            responseWithOrigin(w, r, 200, []byte(err.Error()))
            return
        }
        responseWithOrigin(w, r, 200, b)
    })

    http.HandleFunc("/rank", func(w http.ResponseWriter, r *http.Request) {
        res, err := pkg.Rank(context.TODO(), r.URL.Query().Get("name"))
        if err != nil {
            responseWithOrigin(w, r, 200, []byte(err.Error()))
            return
        }
        b, err := json.Marshal(res)
        if err != nil {
            responseWithOrigin(w, r, 200, []byte(err.Error()))
            return
        }
        responseWithOrigin(w, r, 200, b)
    })

    _ = http.ListenAndServe("127.0.0.1:8089", nil)
}

// avoid cors
func responseWithOrigin(w http.ResponseWriter, r *http.Request, code int, json []byte) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
    w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}