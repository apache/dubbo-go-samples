package config

import (
	_ "github.com/pkg/errors"
)

type Config struct {
	Port    int
	Tracing TracingConfig
}

type TracingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

var cfg *Config = &Config{
	Port: 9555,
}

func Address() int {
	return cfg.Port
}

func Tracing() TracingConfig {
	return cfg.Tracing
}
