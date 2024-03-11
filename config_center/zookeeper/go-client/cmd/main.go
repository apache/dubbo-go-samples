package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	greet "github.com/apache/dubbo-go-samples/config_center/zookeeper/proto"
	"github.com/dubbogo/gost/log/logger"
	"time"
)

const configCenterZKClientConfig = `## set in config center, group is 'dubbogo', dataid is 'dubbo-go-samples-configcenter-zookeeper-client', namespace is default
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

func main() {
	dynamicConfig, err := config.NewConfigCenterConfigBuilder().
		SetProtocol("zookeeper").
		SetAddress("127.0.0.1:2181").
		Build().GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}
	if err = dynamicConfig.PublishConfig("dubbo-go-samples-configcenter-zookeeper-client", "dubbogo", configCenterZKClientConfig); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 10)

	zkOption := config_center.WithZookeeper()
	dataIdOption := config_center.WithDataID("dubbo-go-samples-configcenter-zookeeper-client")
	addressOption := config_center.WithAddress("127.0.0.1:2181")
	groupOption := config_center.WithGroup("dubbogo")
	ins, err := dubbo.NewInstance(
		dubbo.WithConfigCenter(zkOption, dataIdOption, addressOption, groupOption),
	)
	if err != nil {
		panic(err)
	}
	// configure the params that only client layer cares
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		panic(err)
	}

	resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
	if err != nil {
		logger.Error(err)
	}
	logger.Infof("Greet response: %s", resp)
}
