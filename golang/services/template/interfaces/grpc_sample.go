package interfaces

import (
	"context"
	"fmt"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
)

type Server struct {
	pb.UnimplementedGreetingServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s", req.GetName()),
	}, nil
}
