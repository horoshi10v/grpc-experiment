package main

import (
	"context"
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
	"github.com/prometheus/client_golang/prometheus"
)

// Реалізація методу RequestResponse для структури server
func (s *server) RequestResponse(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	GrpcRequestsTotal.Inc()                           // Використання глобальної метрики
	timer := prometheus.NewTimer(GrpcRequestDuration) // Використання глобальної метрики
	defer timer.ObserveDuration()

	GrpcActiveRequests.Inc()       // Збільшуємо кількість активних запитів
	defer GrpcActiveRequests.Dec() // Зменшуємо кількість активних запитів після завершення

	log.Printf("Received request: %s", req.Message)
	time.Sleep(100 * time.Millisecond) // Імітація деякої затримки обробки

	return &pb.Response{Message: "Hello " + req.Message}, nil
}
