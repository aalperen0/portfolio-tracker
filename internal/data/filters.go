package data

import (
	"strings"

	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

type Filters struct {
	Ids      string
	Page     int
	PerPage  int
	Order    string
	Sort     string
	SortList []string
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

func ValidateOtherFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page < 100, "page", "must be less than 100")
	v.Check(f.PerPage > 0, "per_page", "must be greater than zero")
	v.Check(f.PerPage < 250, "per_page", "must be less than 250")

	v.Check(validator.PermittedValues(f.Sort, f.SortList...), "sort", "invalid sort value")
}

func (f Filters) Limit() int {
	return f.PerPage
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PerPage
}

func (f Filters) SortColumn() string {
	for _, value := range f.SortList {
		if f.Sort == value {
			if strings.HasSuffix(value, "_desc") {
				return strings.TrimSuffix(f.Sort, "_desc")
			}
			if strings.HasSuffix(value, "_asc") {
				return strings.TrimSuffix(f.Sort, "_asc")
			}
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) SortDirection() string {
	if strings.HasSuffix(f.Sort, "desc") {
		return "DESC"
	}
	return "ASC"
}
