package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
)

func main() {
	fmt.Println("start gRPC client")

	scanner := bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatal("client connection error:", err)
	}
	defer conn.Close()

	client := pb.NewGreetingServiceClient(conn)

	for {
		fmt.Println("1: send Request")
		fmt.Println("2: exit")
		fmt.Print("please enter -> ")

		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			fmt.Println("Please enter your name")
			scanner.Scan()
			name := scanner.Text()
			req := &pb.HelloRequest{Name: name}
			res, err := client.Hello(context.Background(), req)
			if err != nil {
				fmt.Println("client request error:", err)
			} else {
				fmt.Println("server response:", res.GetMessage())
			}
		case "2":
			fmt.Println("exit")
			goto M
		}
	}
M:
}
