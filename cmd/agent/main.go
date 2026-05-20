package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"go-metrics/internal/agent"
)

func main() {
	serverAddress := flag.String("a", "localhost:8080", "HTTP server address")
	reportInterval := flag.Int("r", 10, "report interval in seconds")
	pollInterval := flag.Int("p", 2, "poll interval in seconds")

	flag.Parse()

	metricAgent := agent.NewAgent(
		normalizeServerAddress(*serverAddress),
		time.Duration(*pollInterval)*time.Second,
		time.Duration(*reportInterval)*time.Second,
	)

	log.Printf(
		"agent started: server=%s pollInterval=%ds reportInterval=%ds",
		normalizeServerAddress(*serverAddress),
		*pollInterval,
		*reportInterval,
	)

	metricAgent.Run()
}

func normalizeServerAddress(address string) string {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return address
	}

	return "http://" + address
}
