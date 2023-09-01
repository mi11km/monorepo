package interfaces

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedPingServiceServer
}

func NewPingServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	//return &pb.PingResponse{
	//	Message:   fmt.Sprintf("Echo: %s", req.GetMessage()),
	//	Timestamp: timestamppb.New(time.Now()),
	//}, nil
	stat := status.New(codes.Unknown, "unknown error occurred")
	stat, err := stat.WithDetails(&pb.ErrDetail{
		Detail: "detail message",
	})
	if err != nil {
		return nil, err
	}
	return nil, stat.Err()
}

func (s *Server) PingServerStream(req *pb.PingRequest, steam pb.PingService_PingServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := steam.Send(&pb.PingResponse{
			Message:   fmt.Sprintf("[%d] Echo: %s", i, req.GetMessage()),
			Timestamp: timestamppb.New(time.Now()),
		}); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func (s *Server) PingClientStream(stream pb.PingService_PingClientStreamServer) error {
	msgList := make([]string, 0)
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&pb.PingResponse{
				Message:   fmt.Sprintf("Echo: %v", msgList),
				Timestamp: timestamppb.New(time.Now()),
			})
		}
		if err != nil {
			return err
		}
		msgList = append(msgList, req.GetMessage())
	}
}

func (s *Server) PingBidirectionalStream(stream pb.PingService_PingBidirectionalStreamServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := stream.Send(&pb.PingResponse{
			Message:   fmt.Sprintf("Echo: %s", req.GetMessage()),
			Timestamp: timestamppb.New(time.Now()),
		}); err != nil {
			return err
		}
	}
}
