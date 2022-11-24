package config

import "github.com/caarlos0/env/v6"

// HTTPConfig defines address for resty client connection
type HTTPConfig struct {
	Address string `env:"ADDRESS" envDefault:"https://rdb.altlinux.org/api/export/branch_binary_packages/"`
}
type AppConfig struct {
	Scope string `env:"SCOPE" envDefault:"all"` // in case set to "diff" then only unique stats will be collected for "releases" only stats on higher releases in branch will be collected, for "all" all metrics will be collected
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
