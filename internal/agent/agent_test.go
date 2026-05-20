package agent

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPollCollectsMetrics(t *testing.T) {
	metricAgent := NewAgent(
		"http://localhost:8080",
		2*time.Second,
		10*time.Second,
	)

	metricAgent.Poll()

	gauges, counters := metricAgent.getMetricsCopy()

	if _, ok := gauges["Alloc"]; !ok {
		t.Fatal("expected Alloc gauge metric")
	}

	if _, ok := gauges["RandomValue"]; !ok {
		t.Fatal("expected RandomValue gauge metric")
	}

	if counters["PollCount"] != 1 {
		t.Fatalf("expected PollCount 1, got %d", counters["PollCount"])
	}
}

func TestReportSendsMetrics(t *testing.T) {
	receivedPaths := make([]string, 0)

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("expected Content-Type text/plain, got %s", r.Header.Get("Content-Type"))
		}

		receivedPaths = append(receivedPaths, r.URL.Path)

		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	metricAgent := NewAgent(
		testServer.URL,
		2*time.Second,
		10*time.Second,
	)

	metricAgent.gauges["Alloc"] = 100
	metricAgent.counters["PollCount"] = 2

	metricAgent.Report()

	if len(receivedPaths) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(receivedPaths))
	}

	hasGauge := false
	hasCounter := false

	for _, path := range receivedPaths {
		if strings.HasPrefix(path, "/update/gauge/Alloc/100") {
			hasGauge = true
		}

		if strings.HasPrefix(path, "/update/counter/PollCount/2") {
			hasCounter = true
		}
	}

	if !hasGauge {
		t.Fatal("expected gauge request")
	}

	if !hasCounter {
		t.Fatal("expected counter request")
	}
}
