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

package handler

import (
	"context"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/cartstore"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"
)

type CartService struct {
	Store cartstore.CartStore
}

func NewCartService() *CartService {
	return &CartService{
		Store: cartstore.NewMemoryCartStore(),
	}
}

func (s *CartService) AddItem(ctx context.Context, in *pb.AddItemRequest) (*pb.Empty, error) {
	err := s.Store.AddItem(ctx, in.UserId, in.Item.ProductId, in.Item.Quantity)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *CartService) GetCart(ctx context.Context, in *pb.GetCartRequest) (*pb.Cart, error) {
	cart, err := s.Store.GetCart(ctx, in.UserId)
	if err != nil {
		return &pb.Cart{}, nil
	}
	out := &pb.Cart{}
	out.UserId = in.UserId
	out.Items = cart.Items
	return out, nil
}

func (s *CartService) EmptyCart(ctx context.Context, in *pb.EmptyCartRequest) (*pb.Empty, error) {
	err := s.Store.EmptyCart(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
