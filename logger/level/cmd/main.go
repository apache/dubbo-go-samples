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
	"context"
	"time"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/protocol"

	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/dubbogo/gost/log/logger"

	log "dubbo.apache.org/dubbo-go/v3/logger"
)

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20000),
		),
		dubbo.WithLogger(
			log.WithLevel("warn"),
			log.WithZap(),
		),
	)
	if err != nil {
		panic(err)
	}
	server, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	go server.Serve()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logger.Info("hello dubbogo this is info log")
			logger.Debug("hello dubbogo this is debug log")
			logger.Warn("hello dubbogo this is warn log")
			time.Sleep(time.Second * 1)
		}
	}
}
