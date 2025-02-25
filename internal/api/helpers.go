package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

type envelope map[string]any

// /readID retrieve the "id" URL parameter from the current request context,
// / then convert to a integer and return it. If the operation isn't successfull
// / it returns 0 and and error
// # Parameters
// @r: The incoming HTTP request
// / # Returns
// / - error: Returns an error if retrieved id is invalid, otherwise returns nil

func (h *Handler) readIDParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	if id == "" {
		return "", errors.New("invalid parameter")
	}
	return id, nil
}

func (h *Handler) writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
	headers http.Header,
) error {
	js, err := json.Marshal(data)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to encode JSON response")
		return err
	}

	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// / readJSON reads and decodes a JSON request body into the provided destination structure.
// / It performs validation checks to ensure the JSON is well-formed and meets the expected structure.
// /
// / # Parameters
// / - w: The HTTP response writer (not used in this function but included for handler context).
// / - r: The incoming HTTP request containing the JSON body.
// / - dest: A pointer to the destination variable where the parsed JSON data should be stored.
// /
// / # Returns
// / - error: Returns an error if JSON parsing fails, otherwise returns nil.
// /
// / # Error Handling
// / - Detects syntax errors and reports the exact character issue.
// / - Reports unexpected end-of-file (EOF) errors when JSON is malformed.
// / - Validates type mismatches and indicates the incorrect field or character offset.
// / - Rejects unknown fields that are not defined in the target structure.
// / - Ensures that the body is not empty.
// / - Ensures that the body contains only a single JSON object.
// / - Panics in case of an invalid unmarshal target.

func (h *Handler) readJSON(w http.ResponseWriter, r *http.Request, dest any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dest)
	if err != nil {

		var (
			syntaxError           *json.SyntaxError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
		)

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf(
				"body contains badly-formed JSON (at character %d)",
				syntaxError.Offset,
			)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeError.Field,
				)
			}
			return fmt.Errorf(
				"body contains incorrect JSON type at (character %d)",
				unmarshalTypeError.Offset,
			)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field: ")
			return fmt.Errorf("body  contains unknown key %s", fieldName)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (h *Handler) readURLstring(qs url.Values, key, defaultValue string) string {
	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}
	return str
}

func (h *Handler) readURLint(
	qs url.Values,
	key string,
	defaultValue int,
	v *validator.Validator,
) int {
	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func calculateFields(c *data.Coin, amount float64, price float64) {
	totalNewCost := amount * price
	c.TotalCost += totalNewCost

	c.Amount += amount
	c.PurchasePriceAverage = c.TotalCost / c.Amount
}

func calculatePNL(c *data.Coin, currentPrice float64) {
	newPNL := (c.Amount * currentPrice) - c.TotalCost
	c.PNL = newPNL
}
