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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/mi11km/workspaces/golang/services/template/interfaces/grpc"
)

/* Interceptors */

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	slog.Info(fmt.Sprintf("[pre] unary client interceptor, gRPC method: %s, req: %v", method, req))
	err := invoker(ctx, method, req, reply, cc, opts...)
	slog.Info(fmt.Sprintf("[post] unary client interceptor, gRPC method: %s, res: %v", method, reply))
	return err
}

func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn,
	method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	slog.Info(fmt.Sprintf("[pre] stream client interceptor, gRPC method: %s", method))
	stream, err := streamer(ctx, desc, cc, method, opts...)
	return &ClientStreamWrapper{stream}, err
}

type ClientStreamWrapper struct {
	grpc.ClientStream
}

func (w *ClientStreamWrapper) SendMsg(m interface{}) error {
	slog.Info(fmt.Sprintf("[pre message] send message: %v", m)) // リクエスト送信前に割り込ませる処理
	return w.ClientStream.SendMsg(m)
}

func (w *ClientStreamWrapper) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if !errors.Is(err, io.EOF) {
		slog.Info(fmt.Sprintf("[post message] recv message: %v", m)) // レスポンス受信後に割り込ませる処理
	}
	return err
}

func (w *ClientStreamWrapper) CloseSend() error {
	err := w.ClientStream.CloseSend()
	slog.Info("[post] stream client interceptor, close send")
	return err
}

/* Main */

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	fmt.Println("start gRPC client")

	address := "localhost:8080"
	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithStreamInterceptor(StreamClientInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("client connection error:", err)
	}
	defer conn.Close()

	client := pb.NewPingServiceClient(conn)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println()
		fmt.Println("1: exit")
		fmt.Println("2: send Request")
		fmt.Println("3: send PingServerStream")
		fmt.Println("4: send PingClientStream")
		fmt.Println("5: send PingBidirectionalStream")
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
			md := metadata.New(map[string]string{"type": "unary", "from": "client"})
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			var headerMD, trailerMD metadata.MD
			res, err := client.Ping(
				ctx, &pb.PingRequest{Message: scanner.Text()}, grpc.Header(&headerMD), grpc.Trailer(&trailerMD))
			if err != nil {
				// change handling depending on status code.
				if stat, ok := status.FromError(err); ok {
					fmt.Printf("code: %s, message: %s, detail: %s\n", stat.Code(), stat.Message(), stat.Details())
				} else {
					fmt.Println("client request error:", err)
				}
			} else {
				fmt.Println("header:", headerMD)
				fmt.Println("trailer:", trailerMD)
				fmt.Println("server response:", res.GetMessage(), res.GetTimestamp())
			}
		case "3":
			fmt.Println("Please enter your message")
			scanner.Scan()
			md := metadata.New(map[string]string{"type": "stream", "from": "client"})
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			stream, err := client.PingServerStream(ctx, &pb.PingRequest{Message: scanner.Text()})
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
			md := metadata.New(map[string]string{"type": "stream", "from": "client"})
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			stream, err := client.PingClientStream(ctx)
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
		case "5":
			md := metadata.New(map[string]string{"type": "stream", "from": "client"})
			ctx := metadata.NewOutgoingContext(context.Background(), md)
			stream, err := client.PingBidirectionalStream(ctx)
			if err != nil {
				fmt.Println("client request error:", err)
				return
			}

			sendNum := 5
			fmt.Printf("Please enter your %d messages\n", sendNum)

			var sendEnd, recvEnd bool
			sendCount := 0
			for !(sendEnd && recvEnd) {
				if !sendEnd {
					scanner.Scan()
					if err := stream.Send(&pb.PingRequest{Message: scanner.Text()}); err != nil {
						fmt.Println("stream send error:", err)
					}
					sendCount++
					if sendCount == sendNum {
						if err := stream.CloseSend(); err != nil {
							fmt.Println("stream close send error:", err)
						}
						sendEnd = true
					}
				}

				var headerMD metadata.MD
				if !recvEnd {
					if headerMD == nil {
						headerMD, err = stream.Header()
						if err != nil {
							fmt.Println("stream header error:", err)
						} else {
							fmt.Println("header:", headerMD)
						}
					}
					res, err := stream.Recv()
					if errors.Is(err, io.EOF) {
						fmt.Println("all the responses have already received.")
						recvEnd = true
						continue
					}
					if err != nil {
						fmt.Println("stream receive error:", err)
					}
					fmt.Println(res)
				}
			}
			trailerMD := stream.Trailer()
			fmt.Println("trailer:", trailerMD)
		default:
			fmt.Println("invalid input")
		}
	}
M:
}
