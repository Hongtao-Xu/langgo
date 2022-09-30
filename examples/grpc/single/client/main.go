package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Hongtao-Xu/langgo/core/rpc"
	"github.com/Hongtao-Xu/langgo/examples/grpc/single/client/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := rpc.NewClient(nil, "127.0.0.1:8000", grpc.WithInsecure())

	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	serverClient := pb.NewServerClient(conn)
	helloResponse, err := serverClient.Hello(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	log.Println(helloResponse, err)
}
