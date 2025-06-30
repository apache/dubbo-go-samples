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

package config

import (
	_ "github.com/pkg/errors"
)

type Config struct {
	Address               string
	Tracing               TracingConfig
	AdService             string
	CartService           string
	CheckoutService       string
	CurrencyService       string
	ProductCatalogService string
	RecommendationService string
	ShippingService       string
}

type TracingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

var cfg *Config = &Config{
	Address:               ":8090",
	AdService:             "adservice",
	CartService:           "cartservice",
	CheckoutService:       "checkoutservice",
	CurrencyService:       "currencyservice",
	ProductCatalogService: "productcatalogservice",
	RecommendationService: "recommendationservice",
	ShippingService:       "shippingservice",
}

func Get() Config {
	return *cfg
}

func Address() string {
	return cfg.Address
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

//func Load() error {
//	configor, err := config.NewConfig(config.WithSource(env.NewSource()))
//	if err != nil {
//		return errors.Wrap(err, "configor.New")
//	}
//	if err := configor.Load(); err != nil {
//		return errors.Wrap(err, "configor.Load")
//	}
//	if err := configor.Scan(cfg); err != nil {
//		return errors.Wrap(err, "configor.Scan")
//	}
//	return nil
//}
