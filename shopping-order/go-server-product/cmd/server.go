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
	"database/sql"
	"fmt"
	"github.com/apache/dubbo-go-samples/shopping-order/go-server-product/pkg"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	"github.com/apache/dubbo-go/common/logger"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/nacos"
	_ "github.com/apache/dubbo-go/registry/protocol"

	"github.com/opentrx/mysql"

	"github.com/transaction-wg/seata-golang/pkg/client"
	seataConfig "github.com/transaction-wg/seata-golang/pkg/client/config"
)

import (
	productDao "github.com/apache/dubbo-go-samples/shopping-order/go-server-product/pkg/dao"
)

const (
	SEATA_CONF_FILE = "SEATA_CONF_FILE"
)

var (
	survivalTimeout = int(3e9)
)

// they are necessary:
// 		export CONF_PROVIDER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
//      export SEATA_CONF_FILE="xxx"
func main() {
	confFile := os.Getenv(SEATA_CONF_FILE)
	seataConfig.InitConf(confFile)
	client.NewRpcClient()

	// init mysql driver of `opentrx/mysql`
	mysql.InitDataResourceManager()
	mysql.RegisterResource(seataConfig.GetATConfig().DSN)

	sqlDB, err := sql.Open("mysql", seataConfig.GetATConfig().DSN)
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(4 * time.Hour)

	d := &productDao.Dao{
		DB: sqlDB,
	}

	svc := &pkg.ProductSvc{
		Dao: d,
	}
	config.SetProviderService(svc)

	config.Load()
	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
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
			fmt.Println("provider app exit now...")
			return
		}
	}
}
