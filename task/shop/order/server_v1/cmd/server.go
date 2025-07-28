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
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/dubbogo/gost/log/logger"
)

import (
	detailAPI "github.com/apache/dubbo-go-samples/task/shop/detail/api"
	"github.com/apache/dubbo-go-samples/task/shop/order/api"
)

// OrderProvider is the provider of order service
type OrderProvider struct {
	detailService detailAPI.Detail
}

func (o *OrderProvider) SubmitOrder(ctx context.Context, req *api.OrderReq) (*api.OrderResp, error) {
	o.detailService.DeductStock(context.Background(), &detailAPI.DeductStockReq{
		Sku:   req.Sku,
		Count: req.Count,
	})
	return &api.OrderResp{
		Env:      "v1",
		Address:  req.Address,
		Phone:    req.Phone,
		Receiver: req.Receiver,
	}, nil
}

func main() {
	ins, err := dubbo.NewInstance(
		dubbo.WithName("shop-order"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20012),
		),
	)
	if err != nil {
		panic(err)
	}

	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}
	detailService, err := detailAPI.NewDetail(cli)
	if err != nil {
		panic(err)
	}
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	if err = api.RegisterOrderHandler(srv, &OrderProvider{detailService: detailService}); err != nil {
		panic(err)
	}
	if err = srv.Serve(); err != nil {
		logger.Error(err)
	}
}
