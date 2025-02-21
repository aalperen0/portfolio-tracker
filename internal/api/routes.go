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
	router.HandlerFunc(http.MethodPatch, "/v1/users/activate", h.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/auth", h.authenticationHandler)
	router.HandlerFunc(http.MethodGet, "/v1/coins", h.GetCoinsFromMarketHandler)
	return h.recoverPanic(h.authenticate(router))
}
