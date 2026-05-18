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