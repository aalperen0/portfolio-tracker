package config

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// / Initalize database connection, pinging db to check connection
// / is established, if after 5 seconds still not established, cancel
// / process.
// / # Return
// / - db: Returns db driver
// / - error: Returns error if connection couldnt established.

func InitDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.DB.maxIdleConns)
	db.SetMaxOpenConns(cfg.DB.maxOpenConns)
	duration, err := time.ParseDuration(cfg.DB.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
