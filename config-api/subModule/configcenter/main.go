package main

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

// todo
func main() {
	centerConfig, err := config.NewConfigCenterConfig(
		config.WithConfigCenterProtocol("zookeeper"),
		config.WithConfigCenterAddress("127.0.0.1:2181")).GetDynamicConfiguration()
	if err != nil {
		panic(err)
	}
	if err := centerConfig.PublishConfig("dubbo-go-samples-configcenter-nacos-server", "dubbo", `dubbo:
  registries:
	demoZK:
	  protocol: zookeeper
	  timeout: 3s
	  address: 127.0.0.1:2181
  protocols:
	triple:
	  name: tri
	  port: 20000
  provider:
	registry:
	  - demoZK
	services:
	  greeterImpl:
		protocol: triple
		interface: com.apache.dubbo.sample.basic.IGreeter # must be compatible with grpc or dubbo-java`); err != nil {
		panic(err)
	}
	select {}
}
