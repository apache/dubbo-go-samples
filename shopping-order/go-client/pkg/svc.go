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
	"errors"
)

import (
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/config"
	seataContext "github.com/transaction-wg/seata-golang/pkg/client/context"
	"github.com/transaction-wg/seata-golang/pkg/client/tm"
)

import (
	"github.com/apache/dubbo-go-samples/shopping-order/go-server-common/filter"
	orderDao "github.com/apache/dubbo-go-samples/shopping-order/go-server-order/pkg/dao"
	productDao "github.com/apache/dubbo-go-samples/shopping-order/go-server-product/pkg/dao"
)

type OrderSvc struct {
	CreateSo func(ctx context.Context, reqs []*orderDao.SoMaster, rsp *orderDao.CreateSoResult) error
}

type ProductSvc struct {
	AllocateInventory func(ctx context.Context, reqs []*productDao.AllocateInventoryReq, rsp *productDao.AllocateInventoryResult) error
}

func (svc *OrderSvc) Reference() string {
	return "OrderSvc"
}

func (svc *ProductSvc) Reference() string {
	return "ProductSvc"
}

var (
	orderSvc   = new(OrderSvc)
	productSvc = new(ProductSvc)
)

func init() {
	config.SetConsumerService(orderSvc)
	config.SetConsumerService(productSvc)
}

type Svc struct {
}

func (svc *Svc) CreateSo(ctx context.Context, rollback bool) ([]uint64, error) {
	rootContext := ctx.(*seataContext.RootContext)
	soMasters := []*orderDao.SoMaster{
		{
			BuyerUserSysNo:       10001,
			SellerCompanyCode:    "SC001",
			ReceiveDivisionSysNo: 110105,
			ReceiveAddress:       "朝阳区长安街001号",
			ReceiveZip:           "000001",
			ReceiveContact:       "Hel",
			ReceiveContactPhone:  "19999999999",
			StockSysNo:           1,
			PaymentType:          1,
			SoAmt:                430.5,
			Status:               10,
			AppId:                "dk-order",
			SoItems: []*orderDao.SoItem{
				{
					ProductSysNo:  1,
					ProductName:   "北冰洋",
					CostPrice:     200,
					OriginalPrice: 232,
					DealPrice:     215.25,
					Quantity:      2,
				},
			},
		},
	}

	reqs := []*productDao.AllocateInventoryReq{{
		ProductSysNo: 1,
		Qty:          2,
	}}

	var createSoResult = &orderDao.CreateSoResult{}
	var allocateInventoryResult = &productDao.AllocateInventoryResult{}

	err1 := orderSvc.CreateSo(context.WithValue(ctx, constant.AttachmentKey, map[string]interface{}{
		filter.SEATA_XID: rootContext.GetXID(),
	}), soMasters, createSoResult)
	if err1 != nil {
		return nil, err1
	}

	err2 := productSvc.AllocateInventory(context.WithValue(ctx, constant.AttachmentKey, map[string]string{
		filter.SEATA_XID: rootContext.GetXID(),
	}), reqs, allocateInventoryResult)
	if err2 != nil {
		return nil, err2
	}

	if rollback {
		return nil, errors.New("there is a error")
	}
	return createSoResult.SoSysNos, nil
}

var service = &Svc{}

type ProxyService struct {
	*Svc
	CreateSo func(ctx context.Context, rollback bool) ([]uint64, error)
}

var methodTransactionInfo = make(map[string]*tm.TransactionInfo)

func init() {
	methodTransactionInfo["CreateSo"] = &tm.TransactionInfo{
		TimeOut:     60000000,
		Name:        "CreateSo",
		Propagation: tm.REQUIRED,
	}
}

func (svc *ProxyService) GetProxyService() interface{} {
	return svc.Svc
}

func (svc *ProxyService) GetMethodTransactionInfo(methodName string) *tm.TransactionInfo {
	return methodTransactionInfo[methodName]
}

var ProxySvc = &ProxyService{
	Svc: service,
}
