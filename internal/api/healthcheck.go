package api

import (
	"net/http"
)

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.config.Env,
			"version":     h.config.Version,
		},
	}

	err := h.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)

	}
}
