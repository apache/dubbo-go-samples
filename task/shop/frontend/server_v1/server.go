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

package server_v1

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/registry"
)

import (
	detailAPI "github.com/apache/dubbo-go-samples/task/shop/detail/api"
	orderAPI "github.com/apache/dubbo-go-samples/task/shop/order/api"
	userAPI "github.com/apache/dubbo-go-samples/task/shop/user/api"
)

// ShopServiceProvider provides the implementation of ShopService interface
type ShopServiceProvider struct {
	userService   userAPI.UserService
	orderService  orderAPI.Order
	detailService detailAPI.Detail
}

func NewShopServiceProvider() (*ShopServiceProvider, error) {
	// global conception
	// configure global configurations and common modules
	ins, err := dubbo.NewInstance(
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress("127.0.0.1:2181"),
		),
	)
	if err != nil {
		panic(err)
	}
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}
	userService, err := userAPI.NewUserService(cli)
	if err != nil {
		panic(err)
	}
	order, err := orderAPI.NewOrder(cli)
	if err != nil {
		panic(err)
	}
	detail, err := detailAPI.NewDetail(cli)
	if err != nil {
		panic(err)
	}
	sp := &ShopServiceProvider{
		userService:   userService,
		orderService:  order,
		detailService: detail,
	}
	return sp, nil
}

// Register registers a user
func (s *ShopServiceProvider) Register(username, password, realName, mail, phone string) bool {
	user := &userAPI.User{
		Username: username,
		Password: password,
		RealName: realName,
		Mail:     mail,
		Phone:    phone,
	}
	if reply, err := s.userService.Register(context.Background(), user); err != nil || !reply.Success {
		return false
	}
	return true
}

func (s *ShopServiceProvider) Login(username, password string) bool {
	req := &userAPI.LoginReq{
		Username: username,
		Password: password,
	}
	if reply, err := s.userService.Login(context.Background(), req); err != nil || reply == nil {
		return false
	}
	return true
}

func (s *ShopServiceProvider) GetUserInfo(username string) (*userAPI.User, error) {
	req := &userAPI.GetInfoReq{
		Username: username,
	}
	return s.userService.GetInfo(context.Background(), req)
}

func (s *ShopServiceProvider) TimeoutLogin(username, password string) bool {
	req := &userAPI.LoginReq{
		Username: username,
		Password: password,
	}
	if reply, err := s.userService.TimeoutLogin(context.Background(), req); err != nil || reply == nil {
		return false
	}
	return true
}

func (s *ShopServiceProvider) CheckItem(sku int64, username string) (*detailAPI.Item, error) {
	req := &detailAPI.GetItemReq{
		Sku:      sku,
		UserName: username,
	}
	// add tag
	ctx := context.Background()
	atm := map[string]string{
		"dubbo.tag":       "gray",
		"dubbo.force.tag": "true",
	}
	ctx = context.WithValue(ctx, constant.AttachmentKey, atm)
	return s.detailService.GetItem(ctx, req)
}

func (s *ShopServiceProvider) CheckItemGray(sku int64, username string) (*detailAPI.Item, error) {
	req := &detailAPI.GetItemReq{
		Sku:      sku,
		UserName: username,
	}
	return s.detailService.GetItem(context.Background(), req)
}

func (s *ShopServiceProvider) SubmitOrder(sku int64, count int, address, phone, receiver string) (*orderAPI.OrderResp, error) {
	order := &orderAPI.OrderReq{
		Sku:      sku,
		Count:    int32(count),
		Address:  address,
		Phone:    phone,
		Receiver: receiver,
	}
	return s.orderService.SubmitOrder(context.Background(), order)
}
