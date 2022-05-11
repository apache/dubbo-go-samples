package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"

	"helloworld/protobuf/helloworld"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

type GreeterProvider struct {
	helloworld.UnimplementedGreeterServer
}

func main() {
	config.SetProviderService(&GreeterProvider{})

	nacosConfig := config.NewRegistryConfigWithProtocolDefaultPort("nacos")
	nacosConfig.Address = "172.19.0.3:8848"
	rc := config.NewRootConfigBuilder().
		SetProvider(config.NewProviderConfigBuilder().
			AddService("GreeterProvider", config.NewServiceConfigBuilder().Build()).
			Build()).
		AddProtocol("tripleProtocolKey", config.NewProtocolConfigBuilder().
			SetName("tri").
			SetPort("20001").
			Build()).
		AddRegistry("registryKey", nacosConfig).
		Build()

	// start dubbo-go framework with configuration
	if err := config.Load(config.WithRootConfig(rc)); err != nil {
		logger.Infof("init ERR = %s\n", err.Error())
	}

	select {}
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.User, error) {
	logger.Infof("SayHello in %s", in.String())
	helloworld := &helloworld.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}
	logger.Infof("SayHello out %s", helloworld.String())
	return helloworld, nil
}
