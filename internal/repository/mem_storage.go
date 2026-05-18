package repository

import "sync"

type MemStorage struct {
	mu sync.RWMutex

	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (m *MemStorage) SetGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gauges[name] = value
}

func (m *MemStorage) AddCounter(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counters[name] += value
}