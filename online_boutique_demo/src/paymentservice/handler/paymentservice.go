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
	"dubbo.apache.org/dubbo-go/v3/config"
	"errors"
	payment "github.com/apache/dubbo-go-samples/online_boutique_demo/paymentservice/proto"
	"github.com/dubbogo/gost/log/logger"
	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
	"strconv"
)

type PaymentService struct{}

func (s *PaymentService) Charge(ctx context.Context, in *payment.ChargeRequest) (*payment.ChargeResponse, error) {
	card := creditcard.Card{
		Number: in.CreditCard.CreditCardNumber,
		Cvv:    strconv.FormatInt(int64(in.CreditCard.CreditCardCvv), 10),
		Year:   strconv.FormatInt(int64(in.CreditCard.CreditCardExpirationYear), 10),
		Month:  strconv.FormatInt(int64(in.CreditCard.CreditCardExpirationMonth), 10),
	}

	//Verify credit card information
	if err := card.Validate(); err != nil {
		logger.Errorf("Invalid credit card: %v", err)
		return nil, errors.New("invalid credit card")
	}

	// TODO:
	// Only VISA and mastercard is accepted, other card types (AMEX, dinersclub) will
	// throw UnacceptedCreditCard error.

	//if card.Company.Short != "visa" && card.Company.Short != "mastercard" {
	//	err := errors.New("unaccepted credit card type")
	//	logger.Errorf("Unaccept credit card type: %s", card.Company.Short)
	//	return nil, err
	//}

	logger.Infof(`Transaction processed: %s, Amount: %s%d.%d`, in.CreditCard.CreditCardNumber, in.Amount.CurrencyCode, in.Amount.Units, in.Amount.Nanos)

	transactionID := uuid.New().String()
	return &payment.ChargeResponse{
		TransactionId: transactionID,
	}, nil
}

func init() {
	config.SetProviderService(new(PaymentService))
}
