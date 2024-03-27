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

package integration

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/client"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/isolation"
	"github.com/dubbogo/gost/log/logger"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	greet "github.com/apache/dubbo-go-samples/filter/proto/sentinel"
	"github.com/stretchr/testify/assert"
)

type GreetFun func(ctx context.Context, req *greet.GreetRequest, opts ...client.CallOption) (*greet.GreetResponse, error)

type stateChangeTestListener struct {
	OnTransformToOpenChan     chan struct{}
	OnTransformToHalfOpenChan chan struct{}
	OnTransformToClosedChan   chan struct{}
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	s.OnTransformToClosedChan <- struct{}{}
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	s.OnTransformToOpenChan <- struct{}{}
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	s.OnTransformToHalfOpenChan <- struct{}{}
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

func TestSentinelConsumerFilter_Concurrency(t *testing.T) {
	_, err := isolation.LoadRules([]*isolation.Rule{
		{
			Resource:   "dubbo:consumer:greet.SentinelGreetService:::Greet()",
			MetricType: isolation.Concurrency,
			Threshold:  100,
		},
	})
	assert.NoError(t, err)
	pass, block := CallGreetFunConcurrently(greetService.Greet, "hello world", 150, 5)
	assert.True(t, pass <= 505 && pass >= 495)
	assert.True(t, block <= 255 && block >= 245)
}

func TestSentinelProviderFilter_Concurrency(t *testing.T) {
	_, err := isolation.LoadRules([]*isolation.Rule{
		{
			Resource:   "dubbo:consumer:greet.SentinelGreetService:::Greet()",
			MetricType: isolation.Concurrency,
			Threshold:  100,
		},
	})
	assert.NoError(t, err)
	pass, block := CallGreetFunConcurrently(greetService.Greet, "hello world", 150, 5)
	assert.True(t, pass <= 505 && pass >= 495)
	assert.True(t, block <= 255 && block >= 245)
}

func TestSentinelProviderFilter_QPS(t *testing.T) {
	pass, block := CallGreetFunConcurrently(greetService.GreetWithQPSLimit, "hello world", 10, 30)
	assert.True(t, pass <= 130 && pass >= 70)
	assert.True(t, block <= 230 && block >= 170)
}

func TestSentinelConsumerFilter_ErrorCount(t *testing.T) {
	listener := &stateChangeTestListener{}
	listener.OnTransformToOpenChan = make(chan struct{}, 1)
	listener.OnTransformToClosedChan = make(chan struct{}, 1)
	listener.OnTransformToHalfOpenChan = make(chan struct{}, 1)
	circuitbreaker.RegisterStateChangeListeners(listener)
	_, err := circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=0.9s, recoveryTimeout=0.9s, maxErrorCount=50
		{
			Resource:                     "dubbo:consumer:greet.SentinelGreetService:::GreetWithChanceOfError()",
			Strategy:                     circuitbreaker.ErrorCount,
			RetryTimeoutMs:               900,
			MinRequestAmount:             10,
			StatIntervalMs:               900,
			StatSlidingWindowBucketCount: 10,
			Threshold:                    50,
		},
	})
	assert.NoError(t, err)

	pass, block := CallGreetFunConcurrently(greetService.GreetWithChanceOfError, "error", 1, 50)
	assert.True(t, pass == 0)
	assert.True(t, block == 50)

	select {
	case <-time.After(time.Second):
		t.Error()
	case <-listener.OnTransformToOpenChan:
	}
	// wait half open
	time.Sleep(time.Second)
	pass, block = CallGreetFunConcurrently(greetService.GreetWithChanceOfError, "hello", 1, 50)
	assert.True(t, pass > 0)
	assert.True(t, block < 50)
	select {
	case <-time.After(time.Second):
		t.Error()
	case <-listener.OnTransformToHalfOpenChan:
	}
	select {
	case <-time.After(time.Second):
		t.Error()
	case <-listener.OnTransformToClosedChan:
	}
}
