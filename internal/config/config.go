package config

type ClientConfig struct {
	Address string `env:"ADDRESS" envDefault:"https://rdb.altlinux.org/api/export/branch_binary_packages/"`
}

type Config struct {
	ClientConfig
}

func NewConfig() *Config {
	return &Config{}
}
