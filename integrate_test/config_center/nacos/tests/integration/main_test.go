package integration

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_center/nacos/proto"
	"os"
	"testing"
	"time"
)

const configCenterNacosClientConfig = `## set in config center, group is 'dubbo', dataid is 'dubbo-go-samples-configcenter-nacos-client', namespace is default
dubbo:
  registries:
    demoZK:
      protocol: zookeeper
      timeout: 3s
      address: 127.0.0.1:2181
  consumer:
    references:
      GreeterClientImpl:
        protocol: tri
        interface: com.apache.dubbo.sample.basic.IGreeter 
`

var greeterProvider greet.GreetService

func TestMain(m *testing.M) {
	dynamicConfig, err := config.NewConfigCenterConfigBuilder().
		SetProtocol("nacos").
		SetAddress("127.0.0.1:8848").
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}
	if err = dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-nacos-client", "dubbo", configCenterNacosClientConfig); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)

	nacosOption := config_center.WithNacos()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-nacos-client")
	addressOption := config_center.WithAddress("127.0.0.1:8848")
	groupOption := config_center.WithGroup("dubbo")
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(nacosOption, dataIdOption, addressOption, groupOption),
	)
	if err != nil {
		panic(err)
	}
	// configure the params that only client layer cares
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	greeterProvider, err = greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}
	time.Sleep(3 * time.Second)
	os.Exit(m.Run())
}
