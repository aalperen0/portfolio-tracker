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
		maxLifeTime  string
	}
	Smtp struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	Coins struct {
		ApiKey string
	}
	Redis struct {
		Host string
	}
}

func LoadConfig() *Config {
	var cfg Config
	flag.IntVar(&cfg.Port, "port", 8080, "Application Port")
	flag.StringVar(&cfg.Env, "env", "development", "development|staging|production")
	flag.StringVar(&cfg.Version, "version", "1.0.0", "versioning")

	// DB
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	flag.IntVar(&cfg.DB.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.maxIdleConns, "db-max-idle-conns", 10, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.maxLifeTime, "db-max-lifetime", "15m", "PostgreSQL max connection")
	flag.StringVar(&cfg.DB.maxIdleTime, "db-max-idle-time", "15m", "Postgresql max idle time")

	defaultDSN := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbName)

	flag.StringVar(&cfg.DB.dsn, "db-dsn", defaultDSN, "PostgreSQL DSN")

	// MAIL
	flag.StringVar(&cfg.Smtp.Host, "smtp-host", "smtp.mailtrap.io", "SMTP Host")
	flag.IntVar(&cfg.Smtp.Port, "smtp-port", 2525, "SMTP Port")
	flag.StringVar(&cfg.Smtp.Username, "smtp-username", "44cb4e51308dbf", "SMTP Username")
	flag.StringVar(&cfg.Smtp.Password, "smtp-password", "089d64e1c8616a", "SMTP Password")
	flag.StringVar(
		&cfg.Smtp.Sender,
		"smtp-sender",
		"Portfolio-Tracker Team <no-reply@example.com>",
		"SMTP Sender",
	)

	redisHost := os.Getenv("REDIS_HOST")

	flag.StringVar(&cfg.Redis.Host, "redis-host", redisHost, "Redis HOST")

	// COIN API apiKey
	coinApiKey := os.Getenv("COINS_API_KEY")
	flag.StringVar(&cfg.Coins.ApiKey, "coin-key", coinApiKey, "Market data")

	flag.Parse()

	return &cfg
}
