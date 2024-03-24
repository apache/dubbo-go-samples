package main

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	_ "github.com/seata/seata-go/pkg/imports"
	"github.com/seata/seata-go/pkg/rm/tcc"

	"github.com/dubbogo/gost/log/logger"

	"github.com/seata/seata-go/pkg/integration"
	"github.com/seata/seata-go/pkg/tm"

	SeataClient "github.com/apache/dubbo-go-samples/transcation/seata-go/triple/client/seata-client"
	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/proto"
	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/service"
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
	uP, err := tcc.NewTCCServiceProxy(&service.UserProvider{})
	if err != nil {
		logger.Errorf("userProviderProxyis not tcc service")
		return
	}

	if resp, err := SeataClient.Prepare(uP, ctx, svc); err != nil {
		logger.Infof("response prepare: %v", err)
	} else {
		logger.Infof("get resp %#v", resp)
	}
	defer func() {
		err := SeataClient.CommitOrRollback(ctx, err == nil)
		if err != nil {
			logger.Errorf("response commit of rollback: %v", err)
			return
		}
		logger.Info("complete commit or rollback")
	}()
	// business
}
