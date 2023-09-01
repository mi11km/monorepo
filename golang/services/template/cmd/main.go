package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/mi11km/workspaces/golang/services/template/interfaces"
	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func main() {
	// config
	port := os.Getenv("PORT")
	cfg := DBConfig{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     os.Getenv("MYSQL_PORT"),
		Name:     os.Getenv("MYSQL_DATABASE"),
	}

	// init logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// init database
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// migration (MEMO: マイグレーションは別でプロセスでやったほうがいい)
	type Ping struct {
		gorm.Model
		Message string
	}
	if err := db.AutoMigrate(&Ping{}); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	// init gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	server := grpc.NewServer()

	pb.RegisterPingServiceServer(server, interfaces.NewPingServer())

	reflection.Register(server)

	go func() {
		slog.Info(fmt.Sprintf("gRPC server listening on port: %s", port))
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
