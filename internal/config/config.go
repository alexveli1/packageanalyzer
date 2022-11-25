// Package config holds configuration structure and its methods to set up application environment
package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

// HTTPConfig defines address for resty client connection and timeout
type HTTPConfig struct {
	Address string        `env:"ADDRESS" envDefault:"https://rdb.altlinux.org/api/export/branch_binary_packages/"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"180s"`
}

// AppConfig limiting scope of analysis for usecase layer
type AppConfig struct {
	Scope string `env:"SCOPE" envDefault:"all"` // in case set to "diff" then only unique stats will be collected for "releases" only stats on higher releases in branch will be collected, for "all" all metrics will be collected
}

// Config consolidated config object to pass around to layers
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

// parseFlags overrides environment variables and defaults for Config structures
func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Scope, "s", cfg.Scope, "scope of processing  all | diff | releases.\n"+
		"all - branch unique packages and releases differences;\n"+
		"diff - list branch unique packages;\n"+
		"releases - release differences")
	flag.Parse()
}
