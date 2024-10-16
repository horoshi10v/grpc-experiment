package main

import (
	"github.com/horoshi10v/grpc-experiment/client/pubsub"
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
	"google.golang.org/grpc"
)

func main() {
	// Підключення до gRPC сервера
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Створення клієнта для ExperimentService
	client := pb.NewExperimentServiceClient(conn)

	// Тестування методу RequestResponse
	TestRequestResponse(client)

	// Додавання затримки для наочного тестування потокових запитів
	time.Sleep(1 * time.Second)

	// Тестування асинхронних потоків
	TestStream(client)

	// Тестування публікації/підписки
	pubsub.TestPubSub(client)
}
