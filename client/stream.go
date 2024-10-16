package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// Тестування асинхронної потокової взаємодії
func TestStream(client pb.ExperimentServiceClient) {
	// Збільшуємо тайм-аут до 30 секунд
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	stream, err := client.StreamRequestResponse(ctx)
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	// Надсилаємо 5 повідомлень
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.Request{Message: "Stream message " + time.Now().Format(time.RFC3339)}); err != nil {
			log.Fatalf("Error sending message to stream: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Отримуємо відповіді зі стріму
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream response: %v", err)
		}
		log.Printf("Received from stream: %s", response.Message)
	}
}
