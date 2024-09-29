package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Створення метрик
var (
	GrpcRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests processed",
	})

	GrpcRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Duration of gRPC requests in seconds",
		Buckets: prometheus.DefBuckets,
	})

	GrpcActiveRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "grpc_active_requests",
		Help: "Current number of active gRPC requests",
	})
)

// Функція для реєстрації метрик HTTP сервером
func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":2112", nil); err != nil {
			panic(err)
		}
	}()
}
