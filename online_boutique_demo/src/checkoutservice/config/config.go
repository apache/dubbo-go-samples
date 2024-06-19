package config

import (
	"fmt"
)

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

func Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

// Load TODO:
func Load() error {
	return nil
}
