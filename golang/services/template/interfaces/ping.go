package interfaces

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedPingServiceServer
}

func NewPingServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		slog.Info(fmt.Sprintf("metadata: %v", md))
	}

	// header サーバーがクライアントに送る最初のヘッダーフレーム
	headerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
	if err := grpc.SetHeader(ctx, headerMD); err != nil {
		return nil, err
	}
	// trailer サーバーがクライアントに送る最後のヘッダーフレーム
	trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}

	return &pb.PingResponse{
		Message:   fmt.Sprintf("Echo: %s", req.GetMessage()),
		Timestamp: timestamppb.New(time.Now()),
	}, nil
}

func (s *Server) PingServerStream(req *pb.PingRequest, stream pb.PingService_PingServerStreamServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		slog.Info(fmt.Sprintf("metadata: %v", md))
	}
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&pb.PingResponse{
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
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		slog.Info(fmt.Sprintf("metadata: %v", md))
	}
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
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		slog.Info(fmt.Sprintf("metadata: %v", md))
	}

	headerMD := metadata.New(map[string]string{"type": "bidirectional stream", "from": "server", "in": "header"})
	//if err := stream.SendHeader(headerMD); err != nil {  // すぐにヘッダーを送信したいならばこちら
	//	return err
	//}
	if err := stream.SetHeader(headerMD); err != nil {
		return err
	}

	trailerMD := metadata.New(map[string]string{"type": "bidirectional stream", "from": "server", "in": "trailer"})
	stream.SetTrailer(trailerMD)

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
