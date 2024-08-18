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

package main

import (
	"dubbo.apache.org/dubbo-go/v3"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/registry"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/checkoutservice/handler"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/checkoutservice/proto"
	"github.com/dubbogo/gost/log/logger"
	"os"
)

var (
	name    = "checkoutservice"
	version = "1.0.0"
)

func main() {
	regAddr := os.Getenv("DUBBO_REGISTRY_ADDRESS")
	if regAddr == "" {
		regAddr = "127.0.0.1:2181"
	}

	ins, err := dubbo.NewInstance(
		dubbo.WithName(name),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(regAddr),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20002),
		),
	)
	if err != nil {
		panic(err)
	}
	//server
	srv, err := ins.NewServer()
	if err != nil {
		panic(err)
	}
	client, err := ins.NewClient()
	cartService, err := pb.NewCartService(client)
	currencyService, err := pb.NewCurrencyService(client)
	emailService, err := pb.NewEmailService(client)
	paymentService, err := pb.NewPaymentService(client)
	productCatalogService, err := pb.NewProductCatalogService(client)
	shippingService, err := pb.NewShippingService(client)
	if err := pb.RegisterCheckoutServiceHandler(srv, &handler.CheckoutService{
		CartService:           cartService,
		CurrencyService:       currencyService,
		EmailService:          emailService,
		PaymentService:        paymentService,
		ProductCatalogService: productCatalogService,
		ShippingService:       shippingService,
	}); err != nil {
		panic(err)
	}

	if err := srv.Serve(); err != nil {
		logger.Error(err)
	}
}
