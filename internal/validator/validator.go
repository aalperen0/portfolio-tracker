package validator

import (
	"errors"
	"regexp"
)

var EmailRx = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrDuplicateEmail  = errors.New("duplicate email")
	ErrEditConflict    = errors.New("edit conflict")
	ErrInvalidCurrency = errors.New(
		"invalid currency, please use valid currencies like 'usd', 'gbp', 'try'",
	)
)

// / It contains map of validation errors
type Validator struct {
	Errors map[string]string
}

// / Creates a new instance of Validator with
// / empty error map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// / Valid returns true if Errors map
// / doesn't contain any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// / The function adds an error message to the map
// / if a validation check is not ok
// # Parameters:
// - ok: boolean (check value)
// - key, message: string

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// / If a string value matches with a specific
// / regexp pattern
// # Parameters:
// - value: string
// - rx: regular expression
// # Returns:
// - bool: returns true if it matches
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	for v := range permittedValues {
		if value == permittedValues[v] {
			return true
		}
	}
	return false
}
