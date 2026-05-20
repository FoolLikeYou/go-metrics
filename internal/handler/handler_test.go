package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-metrics/internal/service"
)

type fakeMetricService struct {
	err      error
	value    string
	gauges   map[string]float64
	counters map[string]int64
}

func (f *fakeMetricService) UpdateMetric(metricType string, metricName string, metricValue string) error {
	return f.err
}

func (f *fakeMetricService) GetMetricValue(metricType string, metricName string) (string, error) {
	if f.err != nil {
		return "", f.err
	}

	return f.value, nil
}

func (f *fakeMetricService) GetAllMetrics() (map[string]float64, map[string]int64) {
	return f.gauges, f.counters
}

func TestHandlerUpdateSuccess(t *testing.T) {
	handler := NewHandler(&fakeMetricService{})

	request := httptest.NewRequest(
		http.MethodPost,
		"/update/counter/someMetric/527",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

func TestHandlerUpdateInvalidMetricType(t *testing.T) {
	handler := NewHandler(&fakeMetricService{
		err: service.ErrInvalidMetricType,
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/update/unknown/someMetric/527",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestHandlerUpdateInvalidMetricValue(t *testing.T) {
	handler := NewHandler(&fakeMetricService{
		err: service.ErrInvalidMetricValue,
	})

	request := httptest.NewRequest(
		http.MethodPost,
		"/update/counter/someMetric/abc",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

func TestHandlerGetMetricValueSuccess(t *testing.T) {
	handler := NewHandler(&fakeMetricService{
		value: "527",
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/value/counter/someMetric",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	body := responseRecorder.Body.String()
	if body != "527" {
		t.Fatalf("expected body 527, got %s", body)
	}
}

func TestHandlerGetMetricValueNotFound(t *testing.T) {
	handler := NewHandler(&fakeMetricService{
		err: service.ErrMetricNotFound,
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/value/counter/unknownMetric",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}

func TestHandlerGetAllMetrics(t *testing.T) {
	handler := NewHandler(&fakeMetricService{
		gauges: map[string]float64{
			"Alloc": 100,
		},
		counters: map[string]int64{
			"PollCount": 5,
		},
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	body := responseRecorder.Body.String()

	if !strings.Contains(body, "Alloc") {
		t.Fatal("expected response body to contain Alloc")
	}

	if !strings.Contains(body, "PollCount") {
		t.Fatal("expected response body to contain PollCount")
	}
}

func TestHandlerInternalError(t *testing.T) {
	handler := NewHandler(&fakeMetricService{
		err: errors.New("some error"),
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/value/counter/someMetric",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handler.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, responseRecorder.Code)
	}
}
