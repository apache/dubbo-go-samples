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
	"fmt"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/afex/hystrix-go/hystrix"

	_ "github.com/apache/dubbo-go-extensions/filter/hystrix"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/filter/hystrix/proto"
)

func init() {
	// Configure hystrix command for the GreetService.Greet method
	// Resource name format: dubbo:consumer:InterfaceName:group:version:Method(param1,param2)
	// For this example: dubbo:consumer:greet.GreetService:::Greet(*greet.GreetRequest)
	// Note: The actual param type name is determined by dubbo-go's type reflection
	cmdName := "dubbo:consumer:greet.GreetService:::Greet"

	hystrix.ConfigureCommand(cmdName, hystrix.CommandConfig{
		Timeout:                1000, // 1 second timeout
		MaxConcurrentRequests:  10,   // Max 10 concurrent requests
		RequestVolumeThreshold: 5,    // Minimum 5 requests before circuit can trip
		SleepWindow:            5000, // 5 seconds to wait after circuit opens before testing
		ErrorPercentThreshold:  50,   // 50% error rate triggers circuit opening
	})

	logger.Infof("Configured hystrix command: %s", cmdName)
}

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}

	svc, err := greet.NewGreetService(cli, client.WithFilter("hystrix_consumer"))
	if err != nil {
		panic(err)
	}

	// Test 1: Normal requests
	logger.Info("=== Test 1: Sending normal requests ===")
	for i := 1; i <= 3; i++ {
		resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: fmt.Sprintf("request-%d", i)})
		if err != nil {
			logger.Errorf("Request %d failed: %v", i, err)
		} else {
			logger.Infof("Request %d success: %s", i, resp.Greeting)
		}
	}

	// Test 2: Multiple concurrent requests to potentially trigger circuit breaker
	logger.Info("\n=== Test 2: Sending concurrent requests ===")
	for i := 1; i <= 15; i++ {
		go func(idx int) {
			resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: fmt.Sprintf("concurrent-%d", idx)})
			if err != nil {
				logger.Errorf("Concurrent request %d failed: %v", idx, err)
			} else {
				logger.Infof("Concurrent request %d success: %s", idx, resp.Greeting)
			}
		}(i)
	}

	// Wait for concurrent requests to complete
	time.Sleep(3 * time.Second)

	// Test 3: Try requests after circuit might be open
	logger.Info("\n=== Test 3: Sending requests after concurrent test ===")
	for i := 1; i <= 5; i++ {
		resp, err := svc.Greet(context.Background(), &greet.GreetRequest{Name: fmt.Sprintf("after-%d", i)})
		if err != nil {
			logger.Errorf("After-test request %d failed (circuit might be open): %v", i, err)
		} else {
			logger.Infof("After-test request %d success: %s", i, resp.Greeting)
		}
		time.Sleep(500 * time.Millisecond)
	}

	logger.Info("\nAll tests completed!")
	logger.Info("If you see 'circuit open' errors, it means Hystrix successfully triggered the circuit breaker.")
	logger.Info("Wait a few seconds and try again to see the circuit recover.")
}
