package handler

import (
	"context"

	"github.com/go-micro/demo/cartservice/cartstore"
	pb "github.com/go-micro/demo/cartservice/proto"
)

type CartService struct {
	Store cartstore.CartStore
}

func (s *CartService) AddItem(ctx context.Context, in *pb.AddItemRequest, out *pb.Empty) error {
	return s.Store.AddItem(ctx, in.UserId, in.Item.ProductId, in.Item.Quantity)
}

func (s *CartService) GetCart(ctx context.Context, in *pb.GetCartRequest, out *pb.Cart) error {
	cart, err := s.Store.GetCart(ctx, in.UserId)
	if err != nil {
		return err
	}
	out.UserId = in.UserId
	out.Items = cart.Items
	return nil
}

func (s *CartService) EmptyCart(ctx context.Context, in *pb.EmptyCartRequest, out *pb.Empty) error {
	return s.Store.EmptyCart(ctx, in.UserId)
}
