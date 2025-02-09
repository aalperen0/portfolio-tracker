package config

import "flag"

type Config struct {
	Port    int
	Env     string
	Version string
}

func LoadConfig() Config {
	var cfg Config
	flag.IntVar(&cfg.Port, "port", 8080, "Application Port")
	flag.StringVar(&cfg.Env, "env", "development", "development|staging|production")
	flag.StringVar(&cfg.Version, "version", "1.0.0", "versioning")
	flag.Parse()

	return cfg
}
