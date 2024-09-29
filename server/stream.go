package main

import (
	"log"
	"time"

	pb "github.com/horoshi10v/grpc-experiment/proto"
)

// Реалізація методу StreamRequestResponse
func (s *server) StreamRequestResponse(stream pb.ExperimentService_StreamRequestResponseServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Failed to receive stream request: %v", err)
			return err
		}

		log.Printf("Received stream request: %s", req.Message)

		// Відправка відповіді в стрім
		err = stream.Send(&pb.Response{Message: "Hello " + req.Message})
		if err != nil {
			log.Printf("Failed to send stream response: %v", err)
			return err
		}

		time.Sleep(1 * time.Second) // Імітація затримки обробки
	}
}
