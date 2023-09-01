package interfaces

import (
	"context"
	"fmt"
	"time"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedPingServiceServer
}

func NewPingServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message:   fmt.Sprintf("Echo: %s", req.GetMessage()),
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}
