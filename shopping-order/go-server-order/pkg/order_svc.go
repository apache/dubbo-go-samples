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
)

import (
	"github.com/apache/dubbo-go/common/constant"

	"github.com/opentrx/mysql"
)

import (
	"github.com/apache/dubbo-go-samples/shopping-order/go-server-common/filter"
	orderDao "github.com/apache/dubbo-go-samples/shopping-order/go-server-order/pkg/dao"
)

type OrderSvc struct {
	Dao *orderDao.Dao
}

func (svc *OrderSvc) CreateSo(ctx context.Context, reqs []*orderDao.SoMaster) (*orderDao.CreateSoResult, error) {
	attach := ctx.Value(constant.AttachmentKey).(map[string]interface{})
	val := attach[filter.SEATA_XID]
	xid := val.(string)
	// set transaction xid
	soIds, err := svc.Dao.CreateSO(
		context.WithValue(context.Background(), mysql.XID, xid), reqs)

	if err == nil {
		return &orderDao.CreateSoResult{soIds}, nil
	}
	return &orderDao.CreateSoResult{soIds}, err
}

func (svc *OrderSvc) Reference() string {
	return "OrderSvc"
}
