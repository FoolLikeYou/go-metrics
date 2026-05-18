package repository

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
}