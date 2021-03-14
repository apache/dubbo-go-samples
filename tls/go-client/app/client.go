package main

import (
	"context"
	"fmt"
	getty "github.com/apache/dubbo-getty"
	gxlog "github.com/dubbogo/gost/log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

import (
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common/logger"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"

	_ "github.com/apache/dubbo-go/filter/filter_impl"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

var (
	survivalTimeout int = 10e9
)
func init(){
	clientKeyPath, _ := filepath.Abs("../../certs/ca.key")
	caPemPath, _ := filepath.Abs("../../certs/ca.pem")
	config.SetSslEnabled(true)
	config.SetClientTlsConfigBuilder(&getty.ClientTlsConfigBuilder{
		ClientPrivateKeyPath:          clientKeyPath,
		ClientTrustCertCollectionPath: caPemPath,
	})
}
// they are necessary:
// 		export CONF_CONSUMER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
func main() {
	hessian.RegisterPOJO(&User{})
	config.Load()
	time.Sleep(3e9)

	gxlog.CInfo("\n\n\nstart to test dubbo")
	user := &User{}
	for i := 0;i < 10 ;i ++{
		err := userProvider.GetUser(context.TODO(), []interface{}{"A001"}, user)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
	}

	gxlog.CInfo("response result: %v\n", user)
	initSignal()
}
func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP,
		syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				logger.Warnf("app exit now by force...")
				os.Exit(1)
			})
			// The program exits normally or timeout forcibly exits.
			fmt.Println("app exit now...")
			return
		}
	}
}