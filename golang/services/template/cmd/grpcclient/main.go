package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	fmt.Println("start gRPC client")

	scanner := bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("client connection error:", err)
	}
	defer conn.Close()

	client := pb.NewPingServiceClient(conn)

	for {
		fmt.Println("1: exit")
		fmt.Println("2: send Request")
		fmt.Println("3: send PingServerStream")
		fmt.Println("4: send PingClientStream")
		fmt.Print("please enter -> ")

		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			fmt.Println("exit")
			goto M
		case "2":
			fmt.Println("Please enter your message")
			scanner.Scan()
			msg := scanner.Text()
			res, err := client.Ping(context.Background(), &pb.PingRequest{Message: msg})
			if err != nil {
				fmt.Println("client request error:", err)
			} else {
				fmt.Println("server response:", res.GetMessage(), res.GetTimestamp())
			}
		case "3":
			fmt.Println("Please enter your message")
			scanner.Scan()
			msg := scanner.Text()
			stream, err := client.PingServerStream(context.Background(), &pb.PingRequest{Message: msg})
			if err != nil {
				fmt.Println("client request error:", err)
				return
			}
			for {
				res, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					fmt.Println("all the responses have already received.")
					break
				}
				if err != nil {
					fmt.Println("stream receive error:", err)
				}
				fmt.Println(res)
			}
		case "4":
			stream, err := client.PingClientStream(context.Background())
			if err != nil {
				fmt.Println("client request error:", err)
				return
			}

			sendCount := 5
			fmt.Printf("Please enter your %d messages\n", sendCount)
			for i := 0; i < sendCount; i++ {
				scanner.Scan()
				if err := stream.Send(&pb.PingRequest{Message: scanner.Text()}); err != nil {
					fmt.Println("stream send error:", err)
					return
				}
			}
			res, err := stream.CloseAndRecv()
			if err != nil {
				fmt.Println("stream close and receive error:", err)
				return
			}
			fmt.Println(res)
		}
	}
M:
}
