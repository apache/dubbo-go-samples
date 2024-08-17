package config

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
	Port: 3550,
}

func Address() int {
	return cfg.Port
}

func Tracing() TracingConfig {
	return cfg.Tracing
}
