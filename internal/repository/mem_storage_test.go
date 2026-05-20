package repository

import "testing"

func TestSetGauge(t *testing.T) {
	storage := NewMemStorage()

	storage.SetGauge("temperature", 23.5)
	storage.SetGauge("temperature", 99.1)

	value, ok := storage.GetGauge("temperature")
	if !ok {
		t.Fatal("expected gauge metric")
	}

	if value != 99.1 {
		t.Fatalf("expected 99.1, got %f", value)
	}
}

func TestAddCounter(t *testing.T) {
	storage := NewMemStorage()

	storage.AddCounter("requests", 10)
	storage.AddCounter("requests", 5)

	value, ok := storage.GetCounter("requests")
	if !ok {
		t.Fatal("expected counter metric")
	}

	if value != 15 {
		t.Fatalf("expected 15, got %d", value)
	}
}

func TestGetAllMetrics(t *testing.T) {
	storage := NewMemStorage()

	storage.SetGauge("temperature", 23.5)
	storage.AddCounter("requests", 10)

	gauges := storage.GetAllGauges()
	counters := storage.GetAllCounters()

	if gauges["temperature"] != 23.5 {
		t.Fatalf("expected gauge temperature 23.5, got %f", gauges["temperature"])
	}

	if counters["requests"] != 10 {
		t.Fatalf("expected counter requests 10, got %d", counters["requests"])
	}
}
