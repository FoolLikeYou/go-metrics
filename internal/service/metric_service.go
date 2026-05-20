package service

import (
	"errors"
	"strconv"

	models "go-metrics/internal/model"
	"go-metrics/internal/repository"
)

var (
	ErrInvalidMetricType  = errors.New("invalid metric type")
	ErrInvalidMetricValue = errors.New("invalid metric value")
	ErrMetricNotFound     = errors.New("metric not found")
)

type MetricService struct {
	storage repository.Storage
}

func NewMetricService(storage repository.Storage) *MetricService {
	return &MetricService{
		storage: storage,
	}
}

func (s *MetricService) UpdateMetric(metricType string, metricName string, metricValue string) error {
	switch metricType {
	case models.Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			return ErrInvalidMetricValue
		}

		s.storage.SetGauge(metricName, value)
		return nil

	case models.Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			return ErrInvalidMetricValue
		}

		s.storage.AddCounter(metricName, value)
		return nil

	default:
		return ErrInvalidMetricType
	}
}

func (s *MetricService) GetMetricValue(metricType string, metricName string) (string, error) {
	switch metricType {
	case models.Gauge:
		value, ok := s.storage.GetGauge(metricName)
		if !ok {
			return "", ErrMetricNotFound
		}

		return strconv.FormatFloat(value, 'f', -1, 64), nil

	case models.Counter:
		value, ok := s.storage.GetCounter(metricName)
		if !ok {
			return "", ErrMetricNotFound
		}

		return strconv.FormatInt(value, 10), nil

	default:
		return "", ErrInvalidMetricType
	}
}

func (s *MetricService) GetAllMetrics() (map[string]float64, map[string]int64) {
	gauges := s.storage.GetAllGauges()
	counters := s.storage.GetAllCounters()

	return gauges, counters
}
