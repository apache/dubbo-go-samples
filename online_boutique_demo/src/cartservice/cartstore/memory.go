package cartstore

import (
	"context"
	"sync"

	pb "github.com/go-micro/demo/cartservice/proto"
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
