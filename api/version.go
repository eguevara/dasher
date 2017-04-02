package api

import (
	"net/http"
	"time"

	"github.com/eguevara/dasher/common"
)

const (
	serviceVersion = "version"
)

// VersionResponse stores the version handler response.
type VersionResponse struct {
	Version string `json:"version"`
}

type versionHandler struct {
	version string
}

func (h *versionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(begin time.Time) {
		requestCount(serviceVersion)
		requestLatency(serviceVersion, begin)
	}(time.Now())

	response := VersionResponse{
		Version: h.version,
	}

	common.Respond(w, response, nil)
}

// VersionHandler handles http requests for the /version api.
func VersionHandler(version string) http.Handler {
	return &versionHandler{
		version: version,
	}
}
