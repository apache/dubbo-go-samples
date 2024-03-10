package integration

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_center/nacos/proto"
	"os"
	"testing"
	"time"
)

var greeterProvider greet.GreetService

func TestMain(m *testing.M) {
	nacosOption := config_center.WithNacos()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-nacos-client")
	addressOption := config_center.WithAddress("127.0.0.1:8848")
	groupOption := config_center.WithGroup("myGroup")
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
