package main

import (
	"flag"
	"log"
	"net/http"

	"go-metrics/internal/handler"
	"go-metrics/internal/repository"
	"go-metrics/internal/service"
)

func main() {
	serverAddress := flag.String("a", "localhost:8080", "HTTP server address")

	flag.Parse()

	storage := repository.NewMemStorage()
	metricService := service.NewMetricService(storage)
	metricHandler := handler.NewHandler(metricService)

	log.Printf("server started on %s", *serverAddress)

	if err := http.ListenAndServe(*serverAddress, metricHandler); err != nil {
		log.Fatal(err)
	}
}
