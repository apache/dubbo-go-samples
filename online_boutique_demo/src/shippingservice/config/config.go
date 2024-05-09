package config

import (
	"fmt"

	"github.com/pkg/errors"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/source/env"
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
	Port: 50051,
}

func Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Tracing() TracingConfig {
	return cfg.Tracing
}

func Load() error {
	configor, err := config.NewConfig(config.WithSource(env.NewSource()))
	if err != nil {
		return errors.Wrap(err, "configor.New")
	}
	if err := configor.Load(); err != nil {
		return errors.Wrap(err, "configor.Load")
	}
	if err := configor.Scan(cfg); err != nil {
		return errors.Wrap(err, "configor.Scan")
	}
	return nil
}
