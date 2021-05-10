/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster_impl"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance"
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/filter/filter_impl"
	_ "dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
	getty "github.com/apache/dubbo-getty"
	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-samples/tls/go-server/pkg"
)

// survival time
var (
	survivalTimeout = int(3e9)
)

func init() {
	serverPemPath, _ := filepath.Abs("../certs/server.pem")
	serverKeyPath, _ := filepath.Abs("../certs/server.key")
	caPemPath, _ := filepath.Abs("../certs/ca.pem")
	config.SetSslEnabled(true)
	config.SetServerTlsConfigBuilder(&getty.ServerTlsConfigBuilder{
		ServerKeyCertChainPath:        serverPemPath,
		ServerPrivateKeyPath:          serverKeyPath,
		ServerTrustCertCollectionPath: caPemPath,
	})
}

/*
	they are necessary:
		export CONF_PROVIDER_FILE_PATH="xx"
		export APP_LOG_CONF_FILE="xx"
*/

func main() {
	// serializing at run time
	hessian.RegisterPOJO(&pkg.User{})
	// load configuration
	config.Load()
	// elegant ending procedure
	initSignal()
}

// elegant ending procedure
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
