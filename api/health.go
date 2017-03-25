package api

import (
	"net/http"

	"github.com/eguevara/dasher/common"
)

// HealthResponse stores the health response.
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler handles all http requests to the /health api.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "OK",
	}
	common.Respond(w, response, nil)
}
