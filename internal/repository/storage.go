package repository

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)

	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)

	GetAllGauges() map[string]float64
	GetAllCounters() map[string]int64
}
