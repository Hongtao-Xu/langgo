package main

import (
	"langgo"
	"langgo/core/rpc"
	"langgo/examples/grpc/single/server/pb"
	"langgo/examples/grpc/single/server/service/server"
)

const addr = "localhost:8000"

func main() {
	langgo.Run()
	cg := rpc.NewServer(nil)
	cg.Use(rpc.LogUnaryServerInterceptor())
	gs, err := cg.Server()
	if err != nil {
		panic(err)
	}
	pb.RegisterServerServer(gs, server.Server{})
	err = cg.Run(addr)
	if err != nil {
		panic(err)
	}
}
