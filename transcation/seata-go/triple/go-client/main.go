package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"dubbo.apache.org/dubbo-go/v3/client"
	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/proto"
	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/service"
	_ "github.com/seata/seata-go/pkg/imports"
	"github.com/seata/seata-go/pkg/rm/tcc"

	"github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/seata/seata-go/pkg/integration"
	"github.com/seata/seata-go/pkg/tm"
)

func main() {
	integration.UseDubbo()
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
		client.WithClientSerialization(constant.ProtobufSerialization),
		client.WithClientProtocolTriple(),
	)
	if err != nil {
		panic(err)
	}

	svc, err := proto.NewUserProvider(cli)
	if err != nil {
		logger.Error(err)
	}
	test(svc)
}

func test(svc proto.UserProvider) {
	ctx := tm.Begin(context.Background(), "TestTCCServiceBusinerr")
	business(ctx, svc)
	<-make(chan struct{})
}

func business(ctx context.Context, svc proto.UserProvider) {
	logger.Info(tm.GetXID(ctx))
	uP, err := tcc.NewTCCServiceProxy(svc)
	if err != nil {
		logger.Infof("userProviderProxyis not tcc service")
		return
	}
	proxyOfUserProvider := &service.UserProviderServer{TCCServiceProxy: uP}

	if resp, err := proxyOfUserProvider.Prepare(ctx, &proto.PrepareRequest{}); err != nil {
		logger.Infof("response prepare: %v", err)
	} else {
		logger.Infof("get resp %#v", resp)
	}
}
