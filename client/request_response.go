package main

import (
	"context"
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// Тестування синхронного запиту-відповіді
func TestRequestResponse(client pb.ExperimentServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Надсилаємо запит
	response, err := client.RequestResponse(ctx, &pb.Request{Message: "Hello from client"})
	if err != nil {
		log.Fatalf("Error calling RequestResponse: %v", err)
	}

	log.Printf("Response from server: %s", response.Message)
}
