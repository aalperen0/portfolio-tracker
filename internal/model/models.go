package model

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/aalperen0/portfolio-tracker/internal/data"
)

type Models struct {
	User  data.UserModel
	Token data.TokenModel
}

func NewModels(db *sql.DB) (Models, error) {
	return Models{
		User:  data.UserModel{DB: db},
		Token: data.TokenModel{DB: db},
	}, nil
}
