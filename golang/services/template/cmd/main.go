package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/mi11km/workspaces/golang/services/template/config"
	"github.com/mi11km/workspaces/golang/services/template/infrastructures"
	"github.com/mi11km/workspaces/golang/services/template/interfaces"
	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.New()

	mysql, err := infrastructures.NewMySQL(cfg.MySQL.FormatDSN())
	if err != nil {
		Fatal(err)
	}
	if err := mysql.Ping(); err != nil {
		Fatal(err)
	}
	defer func() {
		if err := mysql.Close(); err != nil {
			Fatal(err)
		}
	}()

	// init gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterPingServiceServer(server, interfaces.NewPingServer())
	reflection.Register(server)

	go func() {
		slog.Info(fmt.Sprintf("gRPC server listening on port: %s", cfg.Port))
		if err := server.Serve(listener); err != nil {
			Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Shutting down gRPC server...")
	server.GracefulStop()
}

func Fatal(err error) {
	slog.Error(err.Error())
	os.Exit(1)
}
