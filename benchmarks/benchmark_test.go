package benchmarks

import (
	"context"
	"testing"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
	"google.golang.org/grpc"
)

// BenchmarkRequestResponse бенчмарк для методу RequestResponse
func BenchmarkRequestResponse(b *testing.B) {
	// Підключаємося до gRPC сервера
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		b.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewExperimentServiceClient(conn)

	// Створюємо контекст для кожного тесту
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Вимірюємо продуктивність виклику RequestResponse
		_, err := client.RequestResponse(ctx, &pb.Request{Message: "Benchmark Test"})
		if err != nil {
			b.Fatalf("RequestResponse failed: %v", err)
		}
	}
}

// BenchmarkStreamRequestResponse бенчмарк для асинхронного стріму
func BenchmarkStreamRequestResponse(b *testing.B) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		b.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewExperimentServiceClient(conn)

	// Створюємо контекст для кожного тесту
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		stream, err := client.StreamRequestResponse(ctx)
		if err != nil {
			b.Fatalf("Failed to create stream: %v", err)
		}

		// Вимірюємо продуктивність надсилання 5 повідомлень
		for j := 0; j < 5; j++ {
			if err := stream.Send(&pb.Request{Message: "Benchmark Stream"}); err != nil {
				b.Fatalf("Failed to send stream request: %v", err)
			}
		}

		// Вимірюємо продуктивність отримання відповідей
		for {
			_, err := stream.Recv()
			if err != nil {
				break
			}
		}
	}
}
