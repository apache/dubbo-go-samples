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

package api

import (
	detailAPI "github.com/apache/dubbo-go-samples/task/shop/detail/api"
	orderAPI "github.com/apache/dubbo-go-samples/task/shop/order/api"
	userAPI "github.com/apache/dubbo-go-samples/task/shop/user/api"
)

type ShopService interface {
	Register(username, password, realName, mail, phone string) bool

	Login(username, password string) bool

	GetUserInfo(username string) (*userAPI.User, error)

	TimeoutLogin(username, password string) bool

	CheckItem(sku int64, username string) (*detailAPI.Item, error)

	CheckItemGray(sku int64, username string) (*detailAPI.Item, error)

	SubmitOrder(sku int64, count int, address, phone, receiver string) (*orderAPI.OrderResp, error)
}
