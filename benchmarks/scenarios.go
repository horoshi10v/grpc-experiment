package benchmarks

import (
	"context"
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// HighLoadScenario тестує запити при високому навантаженні
func HighLoadScenario(client pb.ExperimentServiceClient, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for {
		_, err := client.RequestResponse(ctx, &pb.Request{Message: "High Load Test"})
		if err != nil {
			log.Printf("Request failed under load: %v", err)
			return
		}
	}
}

// StreamingScenario тестує роботу асинхронних потоків під навантаженням
func StreamingScenario(client pb.ExperimentServiceClient, duration time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	stream, err := client.StreamRequestResponse(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err := stream.Send(&pb.Request{Message: "Streaming Load Test"}); err != nil {
			log.Printf("Failed to send stream request: %v", err)
			return
		}
	}

	for {
		_, err := stream.Recv()
		if err != nil {
			log.Printf("Stream ended: %v", err)
			break
		}
	}
}
