package data

import "github.com/aalperen0/portfolio-tracker/internal/validator"

type Filters struct {
	Ids     string
	Page    int
	PerPage int
	Order   string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page < 100, "page", "must be less than 100")
	v.Check(f.PerPage > 0, "per_page", "must be greater than zero")
	v.Check(f.PerPage < 250, "per_page", "must be less than 250")

	validOrders := []string{"market_cap_asc", "market_cap_desc", "id_asc", "id_desc"}

	v.Check(
		validator.PermittedValues(f.Order, validOrders...),
		"order",
		"must be a valid order type",
	)
}
