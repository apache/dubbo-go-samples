package main

import (
	"fmt"
	getty "github.com/apache/dubbo-getty"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"

	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

// 生存时间
var (
	survivalTimeout = 	int(3e9)
)

func init(){
	serverPemPath, _ := filepath.Abs("../../certs/server.pem")
	serverKeyPath, _ := filepath.Abs("../../certs/server.key")
	caPemPath, _ := filepath.Abs("../../certs/ca.pem")
	config.SetSslEnabled(true)
	config.SetServerTlsConfigBuilder(&getty.ServerTlsConfigBuilder{
		ServerKeyCertChainPath:        serverPemPath,
		ServerPrivateKeyPath:          serverKeyPath,
		ServerTrustCertCollectionPath: caPemPath,
	})
}

/*
	需要配置环境变量
	export CONF_PROVIDER_FILE_PATH="xx"
	export APP_LOG_CONF_FILE="xx"
 */

func main() {
	// 在运行的时候序列化
	hessian.RegisterPOJO(&User{})
	// 加载配置
	config.Load()
	// 优雅的结束程序
	initSignal()
}

// 优雅的结束程序
func initSignal() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})

			fmt.Println("provider app exit now...")
			return
		}
	}
}