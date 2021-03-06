package api

import (
	"net/http"
	"time"

	"github.com/eguevara/dasher/common"
	"github.com/eguevara/dasher/config"
	"github.com/eguevara/go-realtime"
)

const (
	serviceRealtime = "realtime"
)

// realTimeHandler implements a handler.
type realTimeHandler struct {
	Client *realtime.Client
	Config *config.AppConfig
}

// RealTimeHandler handles all https requests to the /realtime api.
func RealTimeHandler(cfg *config.AppConfig) http.Handler {
	client := common.GetOAuthClient(cfg.AnalyticsOAuth)
	realtimeClient := realtime.NewClient(realtime.WithHTTPClient(client))
	return &realTimeHandler{
		Config: cfg,
		Client: realtimeClient,
	}
}

// RealTimeResponse stores the json api response.
type RealTimeResponse struct {
	Report string `json:"report"`
	Value  string `json:"value"`
}

func (rt *realTimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func(begin time.Time) {
		requestCount(serviceRealtime)
		requestLatency(serviceRealtime, begin)
	}(time.Now())

	ids := r.URL.Query().Get("ids")
	metrics := r.URL.Query().Get("metrics")

	response, err := rt.GetMetrics(ids, metrics)
	if err != nil {
		reportRequestFailed(*r, err)
	}

	common.Respond(w, response, err)
}

func (rt *realTimeHandler) GetMetrics(ids, metrics string) (*RealTimeResponse, error) {
	opts := &realtime.Options{
		IDs:     ids,
		Metrics: metrics,
	}

	resp, err := rt.Client.GetRealTime(opts)
	if err != nil {
		return nil, err
	}

	response := &RealTimeResponse{
		Report: metrics,
		Value:  *resp.TotalsForAllResults.RtActiveUsers,
	}
	return response, nil

}
