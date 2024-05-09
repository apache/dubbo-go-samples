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
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/alibaba/sentinel-golang/util"
	greet "github.com/apache/dubbo-go-samples/filter/proto/sentinel"
	"github.com/dubbogo/gost/log/logger"
	"sync"
	"sync/atomic"
	"time"
)

type GreetFun func(ctx context.Context, req *greet.GreetRequest, opts ...client.CallOption) (*greet.GreetResponse, error)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logger.Infof("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	logger.Infof("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	logger.Infof("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewSentinelGreetService(cli, client.WithFilter(constant.SentinelConsumerFilterKey))
	if err != nil {
		panic(err)
	}
	// Register a state change listener so that we could observe the state change of the internal circuit breaker.
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=1s, recoveryTimeout=1s, maxErrorRatio=40%
		{
			Resource:                     "dubbo:consumer:greet.SentinelGreetService:::GreetWithChanceOfError()",
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               2000,
			MinRequestAmount:             10,
			StatIntervalMs:               1000,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    0.4,
		},
	})
	if err != nil {
		panic(err)
	}

	_, err = isolation.LoadRules([]*isolation.Rule{
		{
			Resource:   greet.SentinelGreetService_ServiceInfo.InterfaceName + "::",
			MetricType: isolation.Concurrency,
			Threshold:  100,
		},
	})
	if err != nil {
		panic(err)
	}

	logger.Info("call svc.Greet concurrently")
	CallGreetFunConcurrently(svc.Greet, "hello world", 150, 5)

	logger.Info("call svc.GreetWithQPSLimit concurrently")
	CallGreetFunConcurrently(svc.GreetWithQPSLimit, "hello world", 10, 30)

	logger.Info("call svc.GreetWithChanceOfError triggers circuit breaker open")
	CallGreetFunConcurrently(svc.GreetWithChanceOfError, "error", 1, 300)
	logger.Info("wait circuit breaker HalfOpen")
	time.Sleep(3 * time.Second)
	CallGreetFunConcurrently(svc.GreetWithChanceOfError, "hello world", 1, 300)
	time.Sleep(10 * time.Second)
}

func CallGreetFunConcurrently(f GreetFun, req string, numberOfConcurrently, frequency int) (pass int64, block int64) {
	wg := sync.WaitGroup{}
	wg.Add(numberOfConcurrently)
	for i := 0; i < numberOfConcurrently; i++ {
		go func() {
			for j := 0; j < frequency; j++ {
				_, err := f(context.Background(), &greet.GreetRequest{Name: req})
				if err == nil {
					atomic.AddInt64(&pass, 1)
					//logger.Info(resp.Greeting)
				} else {
					atomic.AddInt64(&block, 1)
					//logger.Error(err)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	logger.Info("success:", pass, "fail:", block)
	return pass, block
}
