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
	Port: 7000,
}

func Address() int {
	return cfg.Port
}

func Tracing() TracingConfig {
	return cfg.Tracing
}
