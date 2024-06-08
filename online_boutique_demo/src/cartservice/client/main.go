package main

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/client"
	hipstershop "github.com/apache/dubbo-go-samples/online_boutique_demo/cartservice/proto"
	"github.com/dubbogo/gost/log/logger"
)

func main() {
	cli, err := client.NewClient(
		client.WithClientURL("tri://127.0.0.1:20000"),
	)
	if err != nil {
		panic(err)
	}
	svc, err := hipstershop.NewCartService(cli)
	if err != nil {
		panic(err)
	}

	respGetcart, err := svc.GetCart(context.Background(), &hipstershop.GetCartRequest{})
	respEmptycart, err := svc.EmptyCart(context.Background(), &hipstershop.EmptyCartRequest{})
	respAdditem, err := svc.AddItem(context.Background(), &hipstershop.AddItemRequest{})

	if err != nil {
		panic(err)
	}
	logger.Infof("get resp get cart: %v", respGetcart)
	logger.Infof("get resp empty cart: %v", respEmptycart)
	logger.Infof("get resp add item: %v", respAdditem)
}
