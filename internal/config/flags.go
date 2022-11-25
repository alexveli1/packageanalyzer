package config

import "flag"

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Scope, "s", cfg.Scope, "scope of processing  all | diff | releases.\n"+
		"all - branch unique packages and releases differences;\n"+
		"diff - list branch unique packages;\n"+
		"releases - release differences")
	flag.Parse()
}
