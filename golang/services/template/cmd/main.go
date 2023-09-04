package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/mi11km/workspaces/golang/services/template/config"
	"github.com/mi11km/workspaces/golang/services/template/interfaces"
	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.New()

	opt := &slog.HandlerOptions{
		//AddSource: true,
		Level: slog.LevelInfo,
	}
	if cfg.Debug {
		opt.Level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opt)))

	// init gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interfaces.UnaryServerLoggingInterceptor,
			interfaces.UnaryServerLoggingInterceptorSub,
		),
		grpc.StreamInterceptor(interfaces.StreamServerInterceptor),
	)

	pb.RegisterPingServiceServer(server, interfaces.NewPingServer())

	reflection.Register(server)

	go func() {
		slog.Info(fmt.Sprintf("gRPC server listening on port: %s", cfg.Port))
		if err := server.Serve(listener); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Shutting down gRPC server...")
	server.GracefulStop()
}
