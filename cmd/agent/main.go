package main

import (
	"log"
	"time"

	"go-metrics/internal/agent"
)

func main() {
	metricAgent := agent.NewAgent(
		"http://localhost:8080",
		2*time.Second,
		10*time.Second,
	)

	log.Println("agent started")

	metricAgent.Run()
}
