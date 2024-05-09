package handler

import (
	"context"
	"strconv"

	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
	"go-micro.dev/v4/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/go-micro/demo/paymentservice/proto"
)

type PaymentService struct{}

func (s *PaymentService) Charge(ctx context.Context, in *pb.ChargeRequest, out *pb.ChargeResponse) error {
	card := creditcard.Card{
		Number: in.CreditCard.CreditCardNumber,
		Cvv:    strconv.FormatInt(int64(in.CreditCard.CreditCardCvv), 10),
		Year:   strconv.FormatInt(int64(in.CreditCard.CreditCardExpirationYear), 10),
		Month:  strconv.FormatInt(int64(in.CreditCard.CreditCardExpirationMonth), 10),
	}
	if err := card.Validate(); err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// TODO:
	// Only VISA and mastercard is accepted, other card types (AMEX, dinersclub) will
	// throw UnacceptedCreditCard error.

	logger.Infof(`Transaction processed: %s, Amount: %s%d.%d`, in.CreditCard.CreditCardNumber, in.Amount.CurrencyCode, in.Amount.Units, in.Amount.Nanos)

	out.TransactionId = uuid.NewString()
	return nil
}
