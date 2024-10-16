package pubsub

import (
	"context"
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// Публікація подій на сервері
func TestPubSub(client pb.ExperimentServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Публікуємо подію
	event := &pb.Event{
		EventId:   "event123",
		EventData: "Test event data",
	}

	ack, err := client.PublishEvent(ctx, event)
	if err != nil {
		log.Fatalf("Error publishing event: %v", err)
	}

	log.Printf("Event published, ack: %s, status: %s", ack.AckId, ack.Status)
}

// Підписка на події з сервера
func SubscribeToEvents(client pb.ExperimentServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Отримуємо потік подій
	stream, err := client.SubscribeEvents(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("Error subscribing to events: %v", err)
	}

	// Читаємо події
	for {
		event, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving event: %v", err)
		}
		log.Printf("Received event: %s, data: %s", event.EventId, event.EventData)
	}
}
