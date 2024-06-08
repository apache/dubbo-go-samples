package cartstore

import "context"
import pb "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"

type CartStore interface {
	AddItem(ctx context.Context, userID, productID string, quantity int32) error
	EmptyCart(ctx context.Context, userID string) error
	GetCart(ctx context.Context, userID string) (*pb.Cart, error)
}

func NewMemoryCartStore() CartStore {
	return &memoryCartStore{
		carts: make(map[string]map[string]int32),
	}
}
