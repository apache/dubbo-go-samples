package integration

import (
	"os"
	"testing"
	"time"

	"dubbo.apache.org/dubbo-go/v3/config"
	hessian "github.com/apache/dubbo-go-hessian2"

	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"

	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"

	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	// _ "dubbo.apache.org/dubbo-go/v3/metadata/remote/impl"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"

	"github.com/apache/dubbo-go-samples/game/pkg/consumer/gate"
	"github.com/apache/dubbo-go-samples/game/pkg/pojo"
)

var gateProvider = new(gate.BasketballService)

func TestMain(m *testing.M) {
	config.SetConsumerService(gateProvider)

	hessian.RegisterPOJO(&pojo.Result{})
	config.Load()
	time.Sleep(3 * time.Second)

	os.Exit(m.Run())
}
