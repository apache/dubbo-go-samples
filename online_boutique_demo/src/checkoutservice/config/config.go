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

type Config struct {
	Port                  int
	Tracing               TracingConfig
	CartService           string
	CurrencyService       string
	EmailService          string
	PaymentService        string
	ProductCatalogService string
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
	Port:                  5050,
	CartService:           "cartservice",
	CurrencyService:       "currencyservice",
	EmailService:          "emailservice",
	PaymentService:        "paymentservice",
	ProductCatalogService: "productcatalogservice",
	ShippingService:       "shippingservice",
}

func Get() Config {
	return *cfg
}

func Address() int {
	return cfg.Port
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

// Load TODO:
func Load() error {
	return nil
}
