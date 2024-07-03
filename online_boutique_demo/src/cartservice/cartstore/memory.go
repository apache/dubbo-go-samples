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

package cartstore

import (
	"context"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"
	"sync"
)

type memoryCartStore struct {
	sync.RWMutex

	carts map[string]map[string]int32
}

func (s *memoryCartStore) AddItem(ctx context.Context, userID, productID string, quantity int32) error {
	s.Lock()
	defer s.Unlock()

	if cart, ok := s.carts[userID]; ok {
		if currentQuantity, ok := cart[productID]; ok {
			cart[productID] = currentQuantity + quantity
		} else {
			cart[productID] = quantity
		}
		s.carts[userID] = cart
	} else {
		s.carts[userID] = map[string]int32{productID: quantity}
	}
	return nil
}

func (s *memoryCartStore) EmptyCart(ctx context.Context, userID string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.carts, userID)
	return nil
}

func (s *memoryCartStore) GetCart(ctx context.Context, userID string) (*pb.Cart, error) {
	s.RLock()
	defer s.RUnlock()

	if cart, ok := s.carts[userID]; ok {
		items := make([]*pb.CartItem, 0, len(cart))
		for p, q := range cart {
			items = append(items, &pb.CartItem{ProductId: p, Quantity: q})
		}
		return &pb.Cart{UserId: userID, Items: items}, nil
	}
	return &pb.Cart{UserId: userID}, nil
}
