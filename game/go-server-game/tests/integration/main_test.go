package integration

import (
    "os"
    "testing"
    "time"

    hessian "github.com/apache/dubbo-go-hessian2"
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

    "github.com/apache/dubbo-go-samples/game/pkg/consumer/game"
    "github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

var gameProvider = new(game.BasketballService)

func TestMain(m *testing.M) {
    config.SetConsumerService(gameProvider)

    hessian.RegisterPOJO(&pojo.Result{})
    config.Load()
    time.Sleep(3 * time.Second)

    os.Exit(m.Run())
}