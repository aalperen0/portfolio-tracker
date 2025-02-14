package model

import (
	"database/sql"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	_ "github.com/lib/pq"
)

type Models struct {
	User data.UserModel
}

func NewModels(db *sql.DB) (Models, error) {

	return Models{
		User: data.UserModel{DB: db},
	}, nil
}
