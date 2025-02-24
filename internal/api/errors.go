package api

import (
	"fmt"
	"net/http"
)

// / Logging with level ERROR with  request method
// / and request URL
// / # Parameters
// / - r:  The incoming HTTP request
// / - err: error type
func (h *Handler) logError(r *http.Request, err error) {
	h.logger.Error().
		Err(err).
		Str("request_method", r.Method).
		Str("request_url", r.URL.String()).
		Msg("request error")
}

// / Sending error response to the client with status and message
// / Also logging the error
// / The response errors below depends on this function and each of them
// / sending appropriate response to the client in specific circumstances
func (h *Handler) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}
	err := h.writeJSON(w, status, env, nil)
	if err != nil {
		h.logError(r, err)
		w.WriteHeader(500)
	}
}

// / The serverErrorResponse() method will be used  application encounters an
// / unexpected problem at runtime. It logs the detailed error message, then uses the
// / errorResponse() helper to send a 500 Internal Server Error status code and JSON
// / response (containing a generic error message) to the client

func (h *Handler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.logError(r, err)
	msg := "the server encountered a problem and couldnt process your request"
	h.errorResponse(w, r, http.StatusInternalServerError, msg)
}

// / The badRequestResponse used to send a  400 Bad Request status code.
func (h *Handler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// / The notFoundResponse() method will be used to send a 404 Not Found
// / status code and  JSON response to the client.
func (h *Handler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	h.errorResponse(w, r, http.StatusNotFound, msg)
}

// / The methodNotAllowedResponse() method will be used to send a 405 Method Not Allowed
// / status code and JSON response to the client
func (h *Handler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("the %s method is not supported by the resource", r.Method)
	h.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

// / The failedValidationResponse() method will be used to send a 422
// / status code and JSON response to the client due to
// / validation process encountered with a problem.
func (h *Handler) failedValidationResponse(
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	h.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (h *Handler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	msg := "unable to update the record due to an edit conflict, please try again"
	h.errorResponse(w, r, http.StatusConflict, msg)
}

func (h *Handler) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	msg := "invalid or wrong credentials"
	h.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (h *Handler) invalidAuthTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	msg := "invalid or missing authentication token"
	h.errorResponse(w, r, http.StatusUnauthorized, msg)
}

func (h *Handler) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	msg := "you must be authenticated to access this resource"
	h.errorResponse(w, r, http.StatusUnauthorized, msg)
}
