package handler

import (
	"context"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/cartstore"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"
)

type CartService struct {
	Store cartstore.CartStore
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
