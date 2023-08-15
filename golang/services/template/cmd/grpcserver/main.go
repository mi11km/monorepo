package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/mi11km/workspaces/golang/services/template/interfaces"
	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
)

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	pb.RegisterGreetingServiceServer(server, interfaces.NewServer())

	reflection.Register(server)

	go func() {
		log.Printf("gRPC server listening on port %d", port)
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down gRPC server...")
	server.GracefulStop()
}
