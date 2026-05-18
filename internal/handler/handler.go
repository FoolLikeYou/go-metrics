package handler

import (
	"errors"
	"net/http"
	"strings"

	"go-metrics/internal/service"
)

type MetricService interface {
	UpdateMetric(metricType string, metricName string, metricValue string) error
}

type Handler struct {
	metricService MetricService
}

func NewHandler(metricService MetricService) *Handler {
	return &Handler{
		metricService: metricService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/update" {
		http.NotFound(w, r)
		return
	}

	if !strings.HasPrefix(r.URL.Path, "/update/") {
		http.NotFound(w, r)
		return
	}

	h.updateMetric(w, r)
}

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	metricType, metricName, metricValue, statusCode, ok := parseUpdatePath(r.URL.Path)
	if !ok {
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	err := h.metricService.UpdateMetric(metricType, metricName, metricValue)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidMetricType):
			http.Error(w, "invalid metric type", http.StatusBadRequest)
			return

		case errors.Is(err, service.ErrInvalidMetricValue):
			http.Error(w, "invalid metric value", http.StatusBadRequest)
			return

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func parseUpdatePath(path string) (
	metricType string,
	metricName string,
	metricValue string,
	statusCode int,
	ok bool,
) {
	parts := strings.Split(path, "/")

	if len(parts) != 5 {
		return "", "", "", http.StatusNotFound, false
	}

	if parts[1] != "update" {
		return "", "", "", http.StatusNotFound, false
	}

	metricType = parts[2]
	metricName = parts[3]
	metricValue = parts[4]

	if metricName == "" {
		return "", "", "", http.StatusNotFound, false
	}

	if metricValue == "" {
		return "", "", "", http.StatusBadRequest, false
	}

	return metricType, metricName, metricValue, http.StatusOK, true
}