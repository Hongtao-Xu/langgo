package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Hongtao-Xu/langgo/core/rpc"
	"github.com/Hongtao-Xu/langgo/examples/grpc/tls/client/pb"
)

const addr = "localhost:8000"

func main() {

	//client证书，私钥；ca证书
	conn, err := rpc.NewClient(&rpc.Tls{
		Crt:   "examples/grpc/tls/keys/client.crt",
		Key:   "examples/grpc/tls/keys/client.key",
		CACrt: "examples/grpc/tls/keys/ca.crt",
	}, addr)
	if err != nil {
		panic(err)
	}

	ServerClient := pb.NewServerClient(conn)

	helloResponse, err := ServerClient.Hello(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Printf("err: %v", err)
		return
	}
	log.Println(helloResponse, err)
}
