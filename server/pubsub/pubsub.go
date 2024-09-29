package pubsub

import (
	"context"
	"log"
	"sync"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// PubSubServer реалізація для публікації та підписки на події
type PubSubServer struct {
	pb.UnimplementedExperimentServiceServer
	subscribers map[string]chan *pb.Event
	mu          sync.Mutex
}

// Створення нового PubSub сервера
func NewPubSubServer() *PubSubServer {
	return &PubSubServer{
		subscribers: make(map[string]chan *pb.Event),
	}
}

// Реалізація методу PublishEvent
func (s *PubSubServer) PublishEvent(ctx context.Context, event *pb.Event) (*pb.EventAck, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Публікація події всім підписникам
	for id, ch := range s.subscribers {
		select {
		case ch <- event:
			log.Printf("Published event to subscriber %s: %s", id, event.EventData)
		default:
			log.Printf("Failed to publish event to subscriber %s: %s", id, event.EventData)
		}
	}

	return &pb.EventAck{AckId: event.EventId, Status: "SUCCESS"}, nil
}

// Реалізація методу SubscribeEvents
func (s *PubSubServer) SubscribeEvents(empty *pb.Empty, stream pb.ExperimentService_SubscribeEventsServer) error {
	s.mu.Lock()
	id := stream.Context().Value("subscriber_id").(string)
	ch := make(chan *pb.Event, 10)
	s.subscribers[id] = ch
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.subscribers, id)
		s.mu.Unlock()
	}()

	for {
		select {
		case event := <-ch:
			if err := stream.Send(event); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}
