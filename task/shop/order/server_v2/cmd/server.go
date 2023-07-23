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

	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	detailAPI "github.com/apache/dubbo-go-samples/task/shop/detail/api"
	"github.com/apache/dubbo-go-samples/task/shop/order/api"
)

// OrderProvider is the provider of order service
type OrderProvider struct {
	api.UnimplementedOrderServer
	detailService *detailAPI.DetailClientImpl
}

func NewOrderProvider() *OrderProvider {
	op := &OrderProvider{}
	// set the detail rpc service
	op.detailService = new(detailAPI.DetailClientImpl)
	config.SetConsumerService(op.detailService)
	return op
}

func (o *OrderProvider) SubmitOrder(ctx context.Context, req *api.OrderReq) (*api.OrderResp, error) {
	o.detailService.DeductStock(context.Background(), &detailAPI.DeductStockReq{
		Sku:   req.Sku,
		Count: req.Count,
	})
	return &api.OrderResp{
		Env:      "v2",
		Address:  req.Address,
		Phone:    req.Phone,
		Receiver: req.Receiver,
	}, nil
}

// export DUBBO_GO_CONFIG_PATH=../conf/dubbogo.yaml
func main() {
	config.SetProviderService(NewOrderProvider())
	if err := config.Load(); err != nil {
		panic(err)
	}
	select {}
}
