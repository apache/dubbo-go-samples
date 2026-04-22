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
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	"github.com/dubbogo/gost/log/logger"
)

import (
	greet "github.com/apache/dubbo-go-samples/graceful_shutdown/proto"
)

func main() {
	addr := flag.String("addr", "tri://127.0.0.1:20000", "server address")
	interval := flag.Duration("interval", 200*time.Millisecond, "interval between requests for each worker")
	shortConn := flag.Bool("short", false, "use short connection (create new client for each request)")
	concurrency := flag.Int("concurrency", 1, "number of concurrent request loops")
	requestTimeout := flag.Duration("request-timeout", 5*time.Second, "per-request timeout")
	namePrefix := flag.String("name-prefix", "hello", "request name prefix")
	maxRequests := flag.Int64("max-requests", 0, "maximum number of requests to issue across all workers, 0 means unlimited")
	minSuccesses := flag.Int64("min-successes", 0, "minimum number of successful requests required before exit")
	minFailures := flag.Int64("min-failures", 0, "minimum number of failed requests required before exit")
	flag.Parse()

	logger.Infof("Starting client, addr=%s short=%v concurrency=%d interval=%s request-timeout=%s",
		*addr, *shortConn, *concurrency, interval.String(), requestTimeout.String())

	var requestCounter atomic.Int64
	var successCount atomic.Int64
	var failureCount atomic.Int64
	var wg sync.WaitGroup
	for workerID := 1; workerID <= *concurrency; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorker(id, *addr, *interval, *shortConn, *requestTimeout, *namePrefix, *maxRequests, &requestCounter, &successCount, &failureCount)
		}(workerID)
	}

	wg.Wait()

	if *maxRequests > 0 {
		successes := successCount.Load()
		failures := failureCount.Load()
		validateRequestSummary(*minSuccesses, successes, *minFailures, failures)
		logger.Infof("Client finished, requests=%d successes=%d failures=%d", requestCounter.Load(), successes, failures)
	}
}

func runWorker(workerID int, addr string, interval time.Duration, shortConn bool, requestTimeout time.Duration, namePrefix string, maxRequests int64, requestCounter, successCount, failureCount *atomic.Int64) {
	var svc greet.GreetService
	var err error

	if !shortConn {
		_, svc, err = newGreetClient(addr)
		if err != nil {
			logger.Errorf("Worker %d failed to create long connection client: %v", workerID, err)
		}
	}

	for {
		requestID, ok := nextRequestID(requestCounter, maxRequests)
		if !ok {
			return
		}

		if shortConn {
			_, svc, err = newGreetClient(addr)
		}
		if err != nil {
			failureCount.Add(1)
			logger.Errorf("Worker %d failed to prepare client: %v", workerID, err)
			time.Sleep(interval)
			continue
		}

		name := fmt.Sprintf("%s-%d", namePrefix, requestID)
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		resp, callErr := svc.Greet(ctx, &greet.GreetRequest{Name: name})
		cancel()

		if callErr != nil {
			failureCount.Add(1)
			logger.Errorf("Worker %d request %d failed: %v", workerID, requestID, callErr)
		} else {
			successCount.Add(1)
			logger.Infof("Worker %d request %d succeeded: %s", workerID, requestID, resp.Greeting)
		}

		time.Sleep(interval)
	}
}

func newGreetClient(addr string) (*client.Client, greet.GreetService, error) {
	cli, err := client.NewClient(client.WithClientURL(addr))
	if err != nil {
		return nil, nil, err
	}

	svc, err := greet.NewGreetService(cli)
	if err != nil {
		return nil, nil, err
	}
	return cli, svc, nil
}

func nextRequestID(counter *atomic.Int64, maxRequests int64) (int64, bool) {
	if maxRequests == 0 {
		return counter.Add(1), true
	}

	for {
		current := counter.Load()
		if current >= maxRequests {
			return 0, false
		}
		next := current + 1
		if counter.CompareAndSwap(current, next) {
			return next, true
		}
	}
}

func validateRequestSummary(minSuccesses, successes, minFailures, failures int64) {
	if successes < minSuccesses {
		panic(fmt.Sprintf("expected at least %d successful requests, got %d", minSuccesses, successes))
	}
	if failures < minFailures {
		panic(fmt.Sprintf("expected at least %d failed requests, got %d", minFailures, failures))
	}
}
