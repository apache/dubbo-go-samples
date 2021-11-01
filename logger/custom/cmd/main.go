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
	"os"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-samples/api"
)

type GreeterProvider struct {
	api.UnimplementedGreeterServer
}

type customLogger struct {
}

func (c *customLogger) Info(args ...interface{}) {
	fmt.Printf("\033[1;32;40m%s\033[0m\n", args)
}

func (c *customLogger) Warn(args ...interface{}) {
	fmt.Printf("\033[1;33;40m%s\033[0m\n", args)
}

func (c *customLogger) Error(args ...interface{}) {
	fmt.Printf("\033[1;31;40m%s\033[0m\n", args)
}

func (c *customLogger) Debug(args ...interface{}) {
	fmt.Printf("\033[1;34;40m%s\033[0m\n", args)
}

func (c *customLogger) Fatal(args ...interface{}) {
	fmt.Printf("\033[1;31;40m%s\033[0m\n", args)
	os.Exit(1)
}

func (c *customLogger) Infof(fmts string, args ...interface{}) {
	fmt.Printf("\033[1;32;40m%s\033[0m\n", fmt.Sprintf(fmts, args...))
}

func (c *customLogger) Warnf(fmts string, args ...interface{}) {
	fmt.Printf("\033[1;33;40m%s\033[0m\n", fmt.Sprintf(fmts, args...))
}

func (c *customLogger) Errorf(fmts string, args ...interface{}) {
	fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf(fmts, args...))
}

func (c *customLogger) Debugf(fmts string, args ...interface{}) {
	fmt.Printf("\033[1;34;40m%s\033[0m\n", fmt.Sprintf(fmts, args...))
}

func (c *customLogger) Fatalf(fmts string, args ...interface{}) {
	fmt.Printf("\033[1;31;40m%s\033[0m\n", fmt.Sprintf(fmts, args...))
	os.Exit(1)
}

func main() {
	config.SetProviderService(&GreeterProvider{})
	if err := config.Load(); err != nil {
		panic(err)
	}
	logger.SetLogger(&customLogger{})
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
