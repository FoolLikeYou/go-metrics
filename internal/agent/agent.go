package agent

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	models "go-metrics/internal/model"
)

type Agent struct {
	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration

	client *http.Client

	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewAgent(serverAddress string, pollInterval time.Duration, reportInterval time.Duration) *Agent {
	return &Agent{
		serverAddress:  serverAddress,
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		client:         &http.Client{Timeout: 5 * time.Second},
		gauges:         make(map[string]float64),
		counters:       make(map[string]int64),
	}
}

func (a *Agent) Run() {
	go func() {
		for {
			a.Poll()
			time.Sleep(a.pollInterval)
		}
	}()

	for {
		a.Report()
		time.Sleep(a.reportInterval)
	}
}

func (a *Agent) Poll() {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	a.mu.Lock()
	defer a.mu.Unlock()

	a.gauges["Alloc"] = float64(memStats.Alloc)
	a.gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	a.gauges["Frees"] = float64(memStats.Frees)
	a.gauges["GCCPUFraction"] = memStats.GCCPUFraction
	a.gauges["GCSys"] = float64(memStats.GCSys)
	a.gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	a.gauges["HeapIdle"] = float64(memStats.HeapIdle)
	a.gauges["HeapInuse"] = float64(memStats.HeapInuse)
	a.gauges["HeapObjects"] = float64(memStats.HeapObjects)
	a.gauges["HeapReleased"] = float64(memStats.HeapReleased)
	a.gauges["HeapSys"] = float64(memStats.HeapSys)
	a.gauges["LastGC"] = float64(memStats.LastGC)
	a.gauges["Lookups"] = float64(memStats.Lookups)
	a.gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	a.gauges["MCacheSys"] = float64(memStats.MCacheSys)
	a.gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	a.gauges["MSpanSys"] = float64(memStats.MSpanSys)
	a.gauges["Mallocs"] = float64(memStats.Mallocs)
	a.gauges["NextGC"] = float64(memStats.NextGC)
	a.gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	a.gauges["NumGC"] = float64(memStats.NumGC)
	a.gauges["OtherSys"] = float64(memStats.OtherSys)
	a.gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	a.gauges["StackInuse"] = float64(memStats.StackInuse)
	a.gauges["StackSys"] = float64(memStats.StackSys)
	a.gauges["Sys"] = float64(memStats.Sys)
	a.gauges["TotalAlloc"] = float64(memStats.TotalAlloc)

	a.gauges["RandomValue"] = rand.Float64()
	a.counters["PollCount"]++
}

func (a *Agent) Report() {
	gauges, counters := a.getMetricsCopy()

	for name, value := range gauges {
		valueString := strconv.FormatFloat(value, 'f', -1, 64)

		err := a.sendMetric(models.Gauge, name, valueString)
		if err != nil {
			fmt.Println(err)
		}
	}

	successfullySentCounters := make([]string, 0, len(counters))

	for name, value := range counters {
		valueString := strconv.FormatInt(value, 10)

		err := a.sendMetric(models.Counter, name, valueString)
		if err != nil {
			fmt.Println(err)
			continue
		}

		successfullySentCounters = append(successfullySentCounters, name)
	}

	a.resetSentCounters(successfullySentCounters)
}

func (a *Agent) getMetricsCopy() (map[string]float64, map[string]int64) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	gauges := make(map[string]float64, len(a.gauges))
	for name, value := range a.gauges {
		gauges[name] = value
	}

	counters := make(map[string]int64, len(a.counters))
	for name, value := range a.counters {
		counters[name] = value
	}

	return gauges, counters
}

func (a *Agent) resetSentCounters(names []string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	for _, name := range names {
		a.counters[name] = 0
	}
}

func (a *Agent) sendMetric(metricType string, metricName string, metricValue string) error {
	url := fmt.Sprintf(
		"%s/update/%s/%s/%s",
		a.serverAddress,
		metricType,
		metricName,
		metricValue,
	)

	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "text/plain")

	response, err := a.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	_, _ = io.Copy(io.Discard, response.Body)

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d for metric %s", response.StatusCode, metricName)
	}

	return nil
}
