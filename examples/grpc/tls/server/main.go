package main

import (
	"github.com/Hongtao-Xu/langgo"
	"github.com/Hongtao-Xu/langgo/core/rpc"
	"github.com/Hongtao-Xu/langgo/examples/grpc/tls/server/pb"
	"github.com/Hongtao-Xu/langgo/examples/grpc/tls/server/service/server"
)

const addr = "localhost:8000"

func main() {
	langgo.Run()
	//client证书，私钥；ca证书
	cg := rpc.NewServer(&rpc.Tls{
		Crt:   "examples/grpc/tls/keys/server.crt",
		Key:   "examples/grpc/tls/keys/server.key",
		CACrt: "examples/grpc/tls/keys/ca.crt",
	})
	cg.Use(rpc.LogUnaryServerInterceptor())
	gs, err := cg.Server()
	if err != nil {
		panic(err)
	}
	pb.RegisterServerServer(gs, server.Server{})
	cg.Run(addr)
}
