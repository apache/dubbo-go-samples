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

package pkg

import (
	"context"
	"fmt"
)

import (
	"github.com/dubbogo/gost/log/logger"

	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
)

import (
	"github.com/apache/dubbo-go-samples/rpc/triple/pb2/api"
	"github.com/apache/dubbo-go-samples/rpc/triple/pb2/models"
)

type GreeterProvider struct {
	api.UnimplementedGreeterServer
}

func (s *GreeterProvider) SayHelloStream(svr api.Greeter_SayHelloStreamServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 1 user, name = %s\n", c.Name)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 2 user, name = %s\n", c2.Name)

	err = svr.Send(&models.User{
		Name: "hello " + c.Name,
		Age:  18,
		ID:   "123456789",
	})
	if err != nil {
		return err
	}
	c3, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go3 GreeterProvider recv 3 user, name = %s\n", c3.Name)

	err = svr.Send(&models.User{
		Name: "hello " + c2.Name,
		Age:  19,
		ID:   "123456789",
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *models.HelloRequest) (*models.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n" + in.Name)
	fmt.Println("get triple header tri-req-id = ", ctx.Value(tripleConstant.TripleCtxKey(tripleConstant.TripleRequestID)))
	fmt.Println("get triple header tri-service-version = ", ctx.Value(tripleConstant.TripleCtxKey(tripleConstant.TripleServiceVersion)))
	return &models.User{Name: "Hello " + in.Name, ID: "12345", Age: 21}, nil
}
