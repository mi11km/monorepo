package interfaces

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"google.golang.org/grpc"
)

func UnaryServerLoggingInterceptor(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info(fmt.Sprintf("[pre] gRPC method: %s", info.FullMethod))
	res, err := handler(ctx, req)
	slog.Info(fmt.Sprintf("[post] gRPC method: %s", info.FullMethod))
	return res, err
}

func UnaryServerLoggingInterceptorSub(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	slog.Info(fmt.Sprintf("[pre sub] gRPC method: %s", info.FullMethod))
	res, err := handler(ctx, req)
	slog.Info(fmt.Sprintf("[post sub] gRPC method: %s", info.FullMethod))
	return res, err
}

func StreamServerInterceptor(srv interface{},
	ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	slog.Info(fmt.Sprintf("[pre stream] open stream, gRPC method: %s", info.FullMethod))
	err := handler(srv, &ServerStreamWrapper{ss})
	slog.Info(fmt.Sprintf("[post stream] close stream, gRPC method: %s", info.FullMethod))
	return err
}

type ServerStreamWrapper struct {
	grpc.ServerStream
}

func (w *ServerStreamWrapper) RecvMsg(m interface{}) error {
	err := w.ServerStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		slog.Info(fmt.Sprintf("[pre message] recv message: %v", m))
	}
	return err
}

func (w *ServerStreamWrapper) SendMsg(m interface{}) error {
	slog.Info(fmt.Sprintf("[post message] send message: %v", m))
	return w.ServerStream.SendMsg(m)
}
