package api

import (
	"errors"
	"net/http"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

func (h *Handler) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return

	}
	err = h.models.User.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			h.failedValidationResponse(w, r, v.Errors)
		default:
			h.serverErrorResponse(w, r, err)
		}
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}

}
