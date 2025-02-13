package model

import (
	"database/sql"
	"fmt"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	_ "github.com/lib/pq"
)

type Models struct {
	User data.UserModel
}

func NewModels(db *sql.DB) (Models, error) {
	if db == nil {
		return Models{}, fmt.Errorf("cannot initialize models: database connection is nil")
	}
	return Models{
		User: data.UserModel{DB: db},
	}, nil
}
