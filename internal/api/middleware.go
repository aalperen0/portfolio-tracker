package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

func (h *Handler) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				h.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			data.ContextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			h.invalidAuthTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateToken(v, token); !v.Valid() {
			h.invalidAuthTokenResponse(w, r)
			return
		}

		user, err := h.models.User.GetUserByToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, validator.ErrRecordNotFound):
				h.invalidAuthTokenResponse(w, r)
			default:
				h.serverErrorResponse(w, r, err)

			}
		}

		r = data.ContextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}
