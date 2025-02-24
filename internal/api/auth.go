package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

func (h *Handler) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlainText(v, input.Password)

	if !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := h.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrRecordNotFound):
			h.invalidCredentialsResponse(w, r)
			return
		default:
			h.serverErrorResponse(w, r, err)
		}
	}

	matches, err := user.Password.Matches(input.Password)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	if !matches {
		h.invalidCredentialsResponse(w, r)
		return
	}

	token, err := h.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
