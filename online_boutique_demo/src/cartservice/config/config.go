package config

import (
	"fmt"
)

type Config struct {
	Port  int
	Redis RedisConfig
}

type RedisConfig struct {
	Addr string
}

var cfg *Config = &Config{
	Port: 7070,
}

func Address() string {
	return fmt.Sprintf(":%d", cfg.Port)
}

func Redis() RedisConfig {
	return cfg.Redis
}

// TODO:Config
func Load() error {
	return nil
}
