package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// / Register the relevant methods, URL patterns and handler functions for
// / endpoints using the HandlerFunc() method.
// / # Return
// / - Returns httprouter instance.
func (h *Handler) Routes() *httprouter.Router {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(h.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(h.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.HealthCheckHandler)
	return router
}
