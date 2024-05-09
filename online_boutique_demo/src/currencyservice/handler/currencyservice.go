package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/go-micro/demo/currencyservice/proto"
)

type CurrencyService struct{}

func (s *CurrencyService) GetSupportedCurrencies(ctx context.Context, in *pb.Empty, out *pb.GetSupportedCurrenciesResponse) error {
	data, err := ioutil.ReadFile("data/currency_conversion.json")
	if err != nil {
		return status.Errorf(codes.Internal, "failed to load currency data: %+v", err)
	}
	currencies := make(map[string]float32)
	if err := json.Unmarshal(data, &currencies); err != nil {
		return status.Errorf(codes.Internal, "failed to unmarshal currency data: %+v", err)
	}
	out.CurrencyCodes = make([]string, 0, len(currencies))
	for k := range currencies {
		out.CurrencyCodes = append(out.CurrencyCodes, k)
	}
	return nil
}

func (s *CurrencyService) Convert(ctx context.Context, in *pb.CurrencyConversionRequest, out *pb.Money) error {
	data, err := ioutil.ReadFile("data/currency_conversion.json")
	if err != nil {
		return status.Errorf(codes.Internal, "failed to load currency data: %+v", err)
	}
	currencies := make(map[string]float64)
	if err := json.Unmarshal(data, &currencies); err != nil {
		return status.Errorf(codes.Internal, "failed to unmarshal currency data: %+v", err)
	}
	fromCurrency, found := currencies[in.From.CurrencyCode]
	if !found {
		return status.Errorf(codes.InvalidArgument, "unsupported currency: %s", in.From.CurrencyCode)
	}
	toCurrency, found := currencies[in.ToCode]
	if !found {
		return status.Errorf(codes.InvalidArgument, "unsupported currency: %s", in.ToCode)
	}
	out.CurrencyCode = in.ToCode
	total := int64(math.Floor(float64(in.From.Units*10^9+int64(in.From.Nanos)) / fromCurrency * toCurrency))
	out.Units = total / 1e9
	out.Nanos = int32(total % 1e9)
	return nil
}
