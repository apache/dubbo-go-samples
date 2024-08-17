package config

type Config struct {
	Port                  int
	Tracing               TracingConfig
	ProductCatalogService string
}

type TracingConfig struct {
	Enable bool
	Jaeger JaegerConfig
}

type JaegerConfig struct {
	URL string
}

var cfg *Config = &Config{
	Port:                  8080,
	ProductCatalogService: "productcatalogservice",
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
