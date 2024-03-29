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
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"github.com/alibaba/sentinel-golang/core/flow"
	greet "github.com/apache/dubbo-go-samples/filter/proto"
	"github.com/dubbogo/gost/log/logger"
	"sync"
	"sync/atomic"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli, client.WithFilter(constant.SentinelConsumerFilterKey))
	if err != nil {
		panic(err)
	}

	// Limit the client's request to GreetService to 200QPS
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               greet.GreetService_ServiceInfo.InterfaceName + "::",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              200,
			RelationStrategy:       flow.CurrentResource,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(10)
	pass := int64(0)
	block := int64(0)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 30; j++ {
				resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: "hello world"})
				if err == nil {
					atomic.AddInt64(&pass, 1)
					logger.Info(resp.Greeting)
				} else {
					atomic.AddInt64(&block, 1)
					logger.Error(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	logger.Info("success:", pass, "fail:", block)
}
