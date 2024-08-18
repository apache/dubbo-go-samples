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
	"fmt"
	"github.com/apache/dubbo-go-samples/online_boutique_demo/frontend/config"
	pb "github.com/apache/dubbo-go-samples/online_boutique_demo/frontend/proto"
	"github.com/dubbogo/gost/log/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

const (
	name    = "frontend"
	version = "1.0.0"

	defaultCurrency = "USD"
	cookieMaxAge    = 60 * 60 * 48

	cookiePrefix    = "shop_"
	cookieSessionID = cookiePrefix + "session-id"
	cookieCurrency  = cookiePrefix + "currency"
)

var (
	whitelistedCurrencies = map[string]bool{
		"USD": true,
		"EUR": true,
		"CAD": true,
		"JPY": true,
		"GBP": true,
		"TRY": true,
	}
)

type ctxKeySessionID struct{}

type frontendServer struct {
	adService             pb.AdService
	cartService           pb.CartService
	checkoutService       pb.CheckoutService
	currencyService       pb.CurrencyService
	productCatalogService pb.ProductCatalogService
	recommendationService pb.RecommendationService
	shippingService       pb.ShippingService
}

func main() {
	regAddr := os.Getenv("DUBBO_REGISTRY_ADDRESS")
	if regAddr == "" {
		regAddr = "127.0.0.1:2181"
	}

	ins, err := dubbo.NewInstance(
		dubbo.WithName("frontendservice"),
		dubbo.WithRegistry(
			registry.WithZookeeper(),
			registry.WithAddress(regAddr),
		),
		dubbo.WithProtocol(
			protocol.WithTriple(),
			protocol.WithPort(20010),
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

	//client
	cli, err := ins.NewClient()
	if err != nil {
		panic(err)
	}

	log := logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout

	//adservice
	adService, err := pb.NewAdService(cli)
	if err != nil {
		panic(err)
	}
	//cartService
	cartService, err := pb.NewCartService(cli)
	if err != nil {
		panic(err)
	}
	//checkoutService
	checkoutService, err := pb.NewCheckoutService(cli)
	if err != nil {
		panic(err)
	}

	//currencyService
	currencyService, err := pb.NewCurrencyService(cli)
	if err != nil {
		panic(err)
	}

	//productcatalog
	productCatalogService, err := pb.NewProductCatalogService(cli)
	if err != nil {
		panic(err)
	}
	//recommendation
	recommendationService, err := pb.NewRecommendationService(cli)

	//shippingService
	shippingService, err := pb.NewShippingService(cli)
	if err != nil {
		panic(err)
	}

	svc := frontendServer{
		adService:             adService,
		cartService:           cartService,
		checkoutService:       checkoutService,
		currencyService:       currencyService,
		productCatalogService: productCatalogService,
		shippingService:       shippingService,
		recommendationService: recommendationService,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/product/{id}", svc.productHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cart", svc.viewCartHandler).Methods(http.MethodGet, http.MethodHead)
	r.HandleFunc("/cart", svc.addToCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/cart/empty", svc.emptyCartHandler).Methods(http.MethodPost)
	r.HandleFunc("/setCurrency", svc.setCurrencyHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", svc.logoutHandler).Methods(http.MethodGet)
	r.HandleFunc("/cart/checkout", svc.placeOrderHandler).Methods(http.MethodPost)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/robots.txt", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "User-agent: *\nDisallow: /") })
	r.HandleFunc("/_healthz", func(w http.ResponseWriter, _ *http.Request) { fmt.Fprint(w, "ok") })

	var handler http.Handler = r
	handler = &logHandler{log: log, next: handler} // add logging
	handler = ensureSessionID(handler)             // add session ID
	// handler = tracing(handler)                     // add opentelemetry instrumentation
	//r.Use(otelmux.Middleware(name))
	//r.Use(tracingContextWrapper)

	//TODO:REGISTER HANDLER
	//if err := micro.RegisterHandler(srv.Server(), handler); err != nil {
	//	logger.Fatal(err)
	//}
	srv_http := &http.Server{
		Addr:    ":8090",
		Handler: handler,
	}

	log.Fatal(srv_http.ListenAndServe())

	logger.Infof("starting server on %s", config.Address())
	if err := srv.Serve(); err != nil {
		logger.Fatal(err)
	}
}
