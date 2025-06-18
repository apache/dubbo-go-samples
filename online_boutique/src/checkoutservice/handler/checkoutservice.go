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
	"fmt"
	"github.com/dubbogo/gost/log/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/apache/dubbo-go-samples/online_boutique_demo/checkoutservice/money"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/checkoutservice/proto"
)

type CheckoutService struct {
	CartService           pb.CartService
	CurrencyService       pb.CurrencyService
	EmailService          pb.EmailService
	PaymentService        pb.PaymentService
	ProductCatalogService pb.ProductCatalogService
	ShippingService       pb.ShippingService
}

func (s *CheckoutService) PlaceOrder(ctx context.Context, in *pb.PlaceOrderRequest) (*pb.PlaceOrderResponse, error) {
	logger.Infof("[PlaceOrder] user_id=%q user_currency=%q", in.UserId, in.UserCurrency)

	orderID, err := uuid.NewUUID()
	if err != nil {
		logger.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to generate order uuid")
	}

	prep, err := s.prepareOrderItemsAndShippingQuoteFromCart(ctx, in.UserId, in.UserCurrency, in.Address)
	if err != nil {
		logger.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	total := &pb.Money{CurrencyCode: in.UserCurrency, Units: 0, Nanos: 0}
	total = money.Must(money.Sum(total, prep.shippingCostLocalized))
	for _, it := range prep.orderItems {
		multPrice := money.MultiplySlow(it.Cost, uint32(it.GetItem().GetQuantity()))
		total = money.Must(money.Sum(total, multPrice))
	}

	txID, err := s.chargeCard(ctx, total, in.CreditCard)
	if err != nil {
		logger.Error(err)
		return nil, status.Errorf(codes.Internal, "failed to charge card: %+v", err)
	}
	logger.Infof("payment went through (transaction_id: %s)", txID)

	shippingTrackingID, err := s.shipOrder(ctx, in.Address, prep.cartItems)
	if err != nil {
		logger.Error(err)
		return nil, status.Errorf(codes.Unavailable, "shipping error: %+v", err)
	}

	if err := s.emptyUserCart(ctx, in.UserId); err != nil {
		logger.Warnf("failed to empty user cart of %s: %+v", in.UserId, err)
	}

	orderResult := &pb.OrderResult{
		OrderId:            orderID.String(),
		ShippingTrackingId: shippingTrackingID,
		ShippingCost:       prep.shippingCostLocalized,
		ShippingAddress:    in.Address,
		Items:              prep.orderItems,
	}

	if err := s.sendOrderConfirmation(ctx, in.Email, orderResult); err != nil {
		logger.Warnf("failed to send order confirmation to %q: %+v", in.Email, err)
	} else {
		logger.Infof("order confirmation email sent to %q", in.Email)
	}
	return &pb.PlaceOrderResponse{Order: orderResult}, nil
}

type orderPrep struct {
	orderItems            []*pb.OrderItem
	cartItems             []*pb.CartItem
	shippingCostLocalized *pb.Money
}

func (s *CheckoutService) prepareOrderItemsAndShippingQuoteFromCart(ctx context.Context, userID, userCurrency string, address *pb.Address) (orderPrep, error) {
	out := orderPrep{}
	cartItems, err := s.getUserCart(ctx, userID)
	if err != nil {
		return out, fmt.Errorf("cart failure: %+v", err)
	}
	orderItems, err := s.prepOrderItems(ctx, cartItems, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to prepare order: %+v", err)
	}
	shippingUSD, err := s.quoteShipping(ctx, address, cartItems)
	if err != nil {
		return out, fmt.Errorf("shipping quote failure: %+v", err)
	}
	shippingPrice, err := s.convertCurrency(ctx, shippingUSD, userCurrency)
	if err != nil {
		return out, fmt.Errorf("failed to convert shipping cost to currency: %+v", err)
	}

	out.shippingCostLocalized = shippingPrice
	out.cartItems = cartItems
	out.orderItems = orderItems
	return out, nil
}

func (s *CheckoutService) quoteShipping(ctx context.Context, address *pb.Address, items []*pb.CartItem) (*pb.Money, error) {
	shippingQuote, err := s.ShippingService.GetQuote(ctx, &pb.GetQuoteRequest{
		Address: address,
		Items:   items,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get shipping quote: %+v", err)
	}
	return shippingQuote.GetCostUsd(), nil
}

func (s *CheckoutService) getUserCart(ctx context.Context, userID string) ([]*pb.CartItem, error) {
	cart, err := s.CartService.GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user cart during checkout: %+v", err)
	}
	return cart.GetItems(), nil
}

func (s *CheckoutService) emptyUserCart(ctx context.Context, userID string) error {
	if _, err := s.CartService.EmptyCart(ctx, &pb.EmptyCartRequest{UserId: userID}); err != nil {
		return fmt.Errorf("failed to empty user cart during checkout: %+v", err)
	}
	return nil
}

func (s *CheckoutService) prepOrderItems(ctx context.Context, items []*pb.CartItem, userCurrency string) ([]*pb.OrderItem, error) {
	out := make([]*pb.OrderItem, len(items))
	for i, item := range items {
		product, err := s.ProductCatalogService.GetProduct(ctx, &pb.GetProductRequest{Id: item.GetProductId()})
		if err != nil {
			return nil, fmt.Errorf("failed to get product #%q", item.GetProductId())
		}
		price, err := s.convertCurrency(ctx, product.GetPriceUsd(), userCurrency)
		if err != nil {
			return nil, fmt.Errorf("failed to convert price of %q to %s", item.GetProductId(), userCurrency)
		}
		out[i] = &pb.OrderItem{Item: item, Cost: price}
	}
	return out, nil
}

func (s *CheckoutService) convertCurrency(ctx context.Context, from *pb.Money, toCurrency string) (*pb.Money, error) {
	result, err := s.CurrencyService.Convert(context.TODO(), &pb.CurrencyConversionRequest{
		From:   from,
		ToCode: toCurrency,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert currency: %+v", err)
	}
	return result, err
}

func (s *CheckoutService) chargeCard(ctx context.Context, amount *pb.Money, paymentInfo *pb.CreditCardInfo) (string, error) {
	paymentResp, err := s.PaymentService.Charge(ctx, &pb.ChargeRequest{
		Amount:     amount,
		CreditCard: paymentInfo,
	})
	if err != nil {
		return "", fmt.Errorf("could not charge the card: %+v", err)
	}
	return paymentResp.GetTransactionId(), nil
}

func (s *CheckoutService) sendOrderConfirmation(ctx context.Context, email string, order *pb.OrderResult) error {
	_, err := s.EmailService.SendOrderConfirmation(ctx, &pb.SendOrderConfirmationRequest{
		Email: email,
		Order: order,
	})
	return err
}

func (s *CheckoutService) shipOrder(ctx context.Context, address *pb.Address, items []*pb.CartItem) (string, error) {
	resp, err := s.ShippingService.ShipOrder(ctx, &pb.ShipOrderRequest{
		Address: address,
		Items:   items,
	})
	if err != nil {
		return "", fmt.Errorf("shipment failed: %+v", err)
	}
	return resp.GetTrackingId(), nil
}
