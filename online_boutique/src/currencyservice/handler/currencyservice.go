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
	"encoding/json"
	"github.com/dubbogo/grpc-go/codes"
	"github.com/dubbogo/grpc-go/status"
	"io/ioutil"
	"math"

	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/currencyservice/proto"
)

type CurrencyService struct{}

func (s *CurrencyService) GetSupportedCurrencies(ctx context.Context, in *pb.Empty) (*pb.GetSupportedCurrenciesResponse, error) {
	out := &pb.GetSupportedCurrenciesResponse{}
	data, err := ioutil.ReadFile("data/currency_conversion.json")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load currency data : %+v", err.Error())
	}
	currencies := make(map[string]float32)
	if err := json.Unmarshal(data, &currencies); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal currency data : %+v", err.Error())
	}
	out.CurrencyCodes = make([]string, 0, len(currencies))
	for k := range currencies {
		out.CurrencyCodes = append(out.CurrencyCodes, k)
	}
	return out, nil
}

func (s *CurrencyService) Convert(ctx context.Context, in *pb.CurrencyConversionRequest) (*pb.Money, error) {
	out := &pb.Money{}
	data, err := ioutil.ReadFile("data/currency_conversion.json")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	currencies := make(map[string]float64)
	if err := json.Unmarshal(data, &currencies); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal currency data: %+v", err)
	}
	fromCurrency, found := currencies[in.From.CurrencyCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency: %s", in.From.CurrencyCode)
	}
	toCurrency, found := currencies[in.ToCode]
	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported currency: %s", in.ToCode)
	}
	out.CurrencyCode = in.ToCode
	total := int64(math.Floor(float64(in.From.Units*10^9+int64(in.From.Nanos)) / fromCurrency * toCurrency))
	out.Units = total / 1e9
	out.Nanos = int32(total % 1e9)
	return out, nil
}
