package tools

import (
	"context" // Видаліть, якщо не потрібно
	"log"
	"net"

	pb "github.com/horoshi10v/grpc-experiment/proto"
	"google.golang.org/grpc"
)

type MockServer struct {
	pb.UnimplementedExperimentServiceServer
}

// RequestResponse імітує відповідь сервера для методу RequestResponse
func (s *MockServer) RequestResponse(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("Mock server received: %s", req.Message)
	return &pb.Response{Message: "Mock Response to " + req.Message}, nil
}

// StreamRequestResponse імітує асинхронний стрім
func (s *MockServer) StreamRequestResponse(stream pb.ExperimentService_StreamRequestResponseServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Mock server stream received: %s", req.Message)
		if err := stream.Send(&pb.Response{Message: "Mock Stream Response"}); err != nil {
			return err
		}
	}
}

// StartMockServer запускає мок-сервер
func StartMockServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExperimentServiceServer(grpcServer, &MockServer{})

	log.Printf("Mock server started on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
