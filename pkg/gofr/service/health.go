package service

import (
	"context"
	"net/http"
)

const (
	serviceUp   = "UP"
	serviceDown = "DOWN"
)

type Health struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details"`
}

func (h *httpService) HealthCheck() *Health {
	var healthResponse = Health{
		Details: make(map[string]interface{}),
	}

	var (
		resp *http.Response
		err  error
	)

	// TODO - Think about a cleaner approach.
	resp, err = h.Get(context.TODO(), h.GetHealthCheckEndpoint(), nil)

	if err != nil || resp == nil {
		healthResponse.Status = serviceDown
		healthResponse.Details["error"] = err.Error()

		return &healthResponse
	}

	defer resp.Body.Close()

	healthResponse.Details["host"] = resp.Request.URL.Host

	if resp.StatusCode == http.StatusOK {
		healthResponse.Status = serviceUp

		return &healthResponse
	}

	healthResponse.Status = serviceDown
	healthResponse.Details["error"] = "service down"

	return &healthResponse
}
