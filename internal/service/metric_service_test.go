package service

import (
	"testing"
)

type fakeStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func newFakeStorage() *fakeStorage {
	return &fakeStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (f *fakeStorage) SetGauge(name string, value float64) {
	f.gauges[name] = value
}

func (f *fakeStorage) AddCounter(name string, value int64) {
	f.counters[name] += value
}

func (f *fakeStorage) GetGauge(name string) (float64, bool) {
	value, ok := f.gauges[name]
	return value, ok
}

func (f *fakeStorage) GetCounter(name string) (int64, bool) {
	value, ok := f.counters[name]
	return value, ok
}

func (f *fakeStorage) GetAllGauges() map[string]float64 {
	return f.gauges
}

func (f *fakeStorage) GetAllCounters() map[string]int64 {
	return f.counters
}

func TestUpdateGauge(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	err := metricService.UpdateMetric("gauge", "temperature", "23.5")
	if err != nil {
		t.Fatal(err)
	}

	if storage.gauges["temperature"] != 23.5 {
		t.Fatalf("expected 23.5, got %f", storage.gauges["temperature"])
	}
}

func TestUpdateCounter(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	err := metricService.UpdateMetric("counter", "requests", "10")
	if err != nil {
		t.Fatal(err)
	}

	err = metricService.UpdateMetric("counter", "requests", "5")
	if err != nil {
		t.Fatal(err)
	}

	if storage.counters["requests"] != 15 {
		t.Fatalf("expected 15, got %d", storage.counters["requests"])
	}
}

func TestGetGaugeValue(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	storage.gauges["temperature"] = 23.5

	value, err := metricService.GetMetricValue("gauge", "temperature")
	if err != nil {
		t.Fatal(err)
	}

	if value != "23.5" {
		t.Fatalf("expected 23.5, got %s", value)
	}
}

func TestGetCounterValue(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	storage.counters["requests"] = 15

	value, err := metricService.GetMetricValue("counter", "requests")
	if err != nil {
		t.Fatal(err)
	}

	if value != "15" {
		t.Fatalf("expected 15, got %s", value)
	}
}

func TestGetUnknownMetric(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	_, err := metricService.GetMetricValue("counter", "unknownMetric")
	if err != ErrMetricNotFound {
		t.Fatalf("expected ErrMetricNotFound, got %v", err)
	}
}

func TestUpdateInvalidMetricType(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	err := metricService.UpdateMetric("unknown", "someMetric", "10")
	if err != ErrInvalidMetricType {
		t.Fatalf("expected ErrInvalidMetricType, got %v", err)
	}
}

func TestUpdateInvalidMetricValue(t *testing.T) {
	storage := newFakeStorage()
	metricService := NewMetricService(storage)

	err := metricService.UpdateMetric("counter", "someMetric", "12.5")
	if err != ErrInvalidMetricValue {
		t.Fatalf("expected ErrInvalidMetricValue, got %v", err)
	}
}
