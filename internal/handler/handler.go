package handler

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go-metrics/internal/service"
)

type MetricService interface {
	UpdateMetric(metricType string, metricName string, metricValue string) error
	GetMetricValue(metricType string, metricName string) (string, error)
	GetAllMetrics() (map[string]float64, map[string]int64)
}

type Handler struct {
	metricService MetricService
	router        chi.Router
}

func NewHandler(metricService MetricService) *Handler {
	h := &Handler{
		metricService: metricService,
		router:        chi.NewRouter(),
	}

	h.initRoutes()

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Handler) initRoutes() {
	h.router.Get("/", h.getAllMetrics)
	h.router.Post("/update/{metricType}/{metricName}/{metricValue}", h.updateMetric)
	h.router.Get("/value/{metricType}/{metricName}", h.getMetricValue)
}

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricName == "" {
		http.NotFound(w, r)
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

func (h *Handler) getMetricValue(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	if metricName == "" {
		http.NotFound(w, r)
		return
	}

	value, err := h.metricService.GetMetricValue(metricType, metricName)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrMetricNotFound):
			http.NotFound(w, r)
			return

		case errors.Is(err, service.ErrInvalidMetricType):
			http.NotFound(w, r)
			return

		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(value))
}

func (h *Handler) getAllMetrics(w http.ResponseWriter, r *http.Request) {
	gauges, counters := h.metricService.GetAllMetrics()

	pageData := struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}{
		Gauges:   gauges,
		Counters: counters,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_ = metricsPageTemplate.Execute(w, pageData)
}

var metricsPageTemplate = template.Must(template.New("metrics").Parse(`
<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<title>Metrics</title>
</head>
<body>
	<h1>Metrics</h1>

	<h2>Gauge</h2>
	<table border="1" cellpadding="5" cellspacing="0">
		<tr>
			<th>Name</th>
			<th>Value</th>
		</tr>
		{{ range $name, $value := .Gauges }}
		<tr>
			<td>{{ $name }}</td>
			<td>{{ $value }}</td>
		</tr>
		{{ else }}
		<tr>
			<td colspan="2">No gauge metrics</td>
		</tr>
		{{ end }}
	</table>

	<h2>Counter</h2>
	<table border="1" cellpadding="5" cellspacing="0">
		<tr>
			<th>Name</th>
			<th>Value</th>
		</tr>
		{{ range $name, $value := .Counters }}
		<tr>
			<td>{{ $name }}</td>
			<td>{{ $value }}</td>
		</tr>
		{{ else }}
		<tr>
			<td colspan="2">No counter metrics</td>
		</tr>
		{{ end }}
	</table>
</body>
</html>
`))
