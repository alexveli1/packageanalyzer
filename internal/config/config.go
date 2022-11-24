package config

import "github.com/caarlos0/env/v6"

type HTTPConfig struct {
	Address string `env:"ADDRESS" envDefault:"https://rdb.altlinux.org/api/export/branch_binary_packages/"`
}
type AppConfig struct {
	Scope string `env:"SCOPE" envDefault:"all"`
}

type Config struct {
	HTTPConfig
	AppConfig
}

func NewConfig() *Config {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {

		return nil
	}
	parseFlags(cfg)

	return cfg
}
