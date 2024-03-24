package main

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/server"

	"github.com/dubbogo/gost/log/logger"

	_ "github.com/seata/seata-go/pkg/imports"
	"github.com/seata/seata-go/pkg/integration"
	"github.com/seata/seata-go/pkg/rm/tcc"

	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/proto"
	"github.com/apache/dubbo-go-samples/transcation/seata-go/triple/service"
)

func main() {
	integration.UseDubbo()
	userProviderProxy, err := tcc.NewTCCServiceProxy(&service.UserProvider{})
	if err != nil {
		logger.Errorf("get userProviderProxy tcc service proxy error, %v", err.Error())
		return
	}
	srv, err := server.NewServer(
		server.WithServerProtocol(
			protocol.WithPort(20000),
			protocol.WithTriple(),
		),
		server.WithServerSerialization(constant.ProtobufSerialization),
	)
	if err != nil {
		panic(err)
	}
	if err := proto.RegisterUserProviderHandler(srv, &service.UserProviderServer{TCCServiceProxy: userProviderProxy}); err != nil {
		panic(err)
	}
	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
