package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// / Register the relevant methods, URL patterns and handler functions for
// / endpoints using the HandlerFunc() method.
// / # Return
// / - Returns httprouter instance.
func (h *Handler) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(h.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.HealthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", h.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/auth", h.authenticationHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/activate", h.activateUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/coins", h.GetCoinsFromMarketHandler)

	protectedRoutes := []struct {
		method  string
		path    string
		handler http.HandlerFunc
	}{
		{http.MethodPost, "/v1/users/coins", h.AddCoinsHandler},
		{http.MethodGet, "/v1/users/coins/:id", h.GetCoinFromPortfolioHandler},
		{http.MethodDelete, "/v1/users/coins/:id", h.DeleteCoinFromPortfolioHandler},
	}

	for _, route := range protectedRoutes {
		router.HandlerFunc(
			route.method,
			route.path,
			h.requiredAuthenticatedUser(http.HandlerFunc(route.handler)),
		)
	}

	return h.recoverPanic(h.authenticate(router))
}
