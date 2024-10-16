package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/horoshi10v/grpc-experiment/proto"
	"github.com/horoshi10v/grpc-experiment/server/pubsub"
	_ "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// Реалізація сервера gRPC
type server struct {
	pb.UnimplementedExperimentServiceServer
	broker *pubsub.KafkaProducer
}

// Реалізація методу PublishEvent з інтеграцією з Kafka
func (s *server) PublishEvent(ctx context.Context, event *pb.Event) (*pb.EventAck, error) {
	// Відстеження активних запитів
	GrpcActiveRequests.Inc()
	defer GrpcActiveRequests.Dec()

	err := s.broker.PublishMessage(event.EventData)
	if err != nil {
		return &pb.EventAck{AckId: event.EventId, Status: "FAILURE"}, err
	}

	// Записуємо загальну кількість запитів
	GrpcRequestsTotal.Inc()

	return &pb.EventAck{AckId: event.EventId, Status: "SUCCESS"}, nil
}

func main() {
	// Ініціалізація метрик (метрики та HTTP сервер для метрик вже реєструються в metrics.go)
	StartMetricsServer()

	// Ініціалізація Kafka Producer
	kafkaProducer := pubsub.NewKafkaProducer()
	defer kafkaProducer.Close()

	// Ініціалізація gRPC сервера з передачею Kafka Producer
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExperimentServiceServer(grpcServer, &server{broker: kafkaProducer})

	// Запуск gRPC сервера у горутині
	go func() {
		log.Printf("Starting gRPC server on port 50051...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Очікування сигналу завершення роботи для коректного завершення
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("Server exited.")
}
