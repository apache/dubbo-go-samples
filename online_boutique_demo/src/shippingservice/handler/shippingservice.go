package handler

import (
	"context"
	"fmt"

	"go-micro.dev/v4/logger"

	pb "github.com/go-micro/demo/shippingservice/proto"
)

type ShippingService struct{}

func (s *ShippingService) GetQuote(ctx context.Context, in *pb.GetQuoteRequest, out *pb.GetQuoteResponse) error {
	logger.Info("[GetQuote] received request")
	defer logger.Info("[GetQuote] completed request")

	// 1. Generate a quote based on the total number of items to be shipped.
	quote := CreateQuoteFromCount(0)

	// 2. Generate a response.
	out.CostUsd = &pb.Money{
		CurrencyCode: "USD",
		Units:        int64(quote.Dollars),
		Nanos:        int32(quote.Cents * 10000000),
	}
	return nil
}

func (s *ShippingService) ShipOrder(ctx context.Context, in *pb.ShipOrderRequest, out *pb.ShipOrderResponse) error {
	logger.Info("[ShipOrder] received request")
	defer logger.Info("[ShipOrder] completed request")
	// 1. Create a Tracking ID
	baseAddress := fmt.Sprintf("%s, %s, %s", in.Address.StreetAddress, in.Address.City, in.Address.State)
	id := CreateTrackingId(baseAddress)

	// 2. Generate a response.
	out.TrackingId = id
	return nil
}
