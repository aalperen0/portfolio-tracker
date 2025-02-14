package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
	"github.com/rs/zerolog"
)

// / Route: POST /v1/users
// / This handler processes incoming HTTP requests for registering a new user.
// / Sends an activation token to activate account.
// / It validates the user's input, hashes the password, checks for errors, and inserts the new user into the database.
// / Request Body: The handler expects a JSON object in the request body with the following fields:
// # Parameters
// @ name (string, required): The name of the user.
// @ email (string, required): The email address of the user.
// @ password (string, required): The password for the user.
// # Response: Success (HTTP Status 201):
// / If the registration is successful, the handler responds with a 201 Created status and the newly created user object.
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

	// Sending email with goroutine at the backround
	// and catching errors before terminating
	go func() {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Level(zerolog.ErrorLevel).With().Err(fmt.Errorf("%s", err))
			}
		}()

		err = h.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			h.logger.Level(zerolog.ErrorLevel).With().Err(err)
		}
	}()

	err = h.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}

}
