package integration

import (
	"os"
	"testing"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
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
