package main

import (
	"log"
	"net/http"

	"go-metrics/internal/handler"
	"go-metrics/internal/repository"
	"go-metrics/internal/service"
)

func main() {
	storage := repository.NewMemStorage()
	metricService := service.NewMetricService(storage)
	metricHandler := handler.NewHandler(metricService)

	if err := http.ListenAndServe(":8080", metricHandler); err != nil {
		log.Fatal(err)
	}
}