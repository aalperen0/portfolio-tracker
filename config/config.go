package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Port    int
	Env     string
	Version string
	DB      struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func LoadConfig() Config {
	var cfg Config
	flag.IntVar(&cfg.Port, "port", 8080, "Application Port")
	flag.StringVar(&cfg.Env, "env", "development", "development|staging|production")
	flag.StringVar(&cfg.Version, "version", "1.0.0", "versioning")

	flag.IntVar(&cfg.DB.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.maxIdleTime, "db-max-idle-time", "15m", "Postgresql max idle time")
	flag.Parse()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	cfg.DB.dsn = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", dbUser, dbPassword, dbPort, dbName)

	return cfg
}
