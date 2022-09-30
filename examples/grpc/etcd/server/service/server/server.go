package server

import (
	"context"
	"github.com/Hongtao-Xu/langgo/examples/grpc/etcd/server/pb"
)

func (s Server) Hello(ctx context.Context, empty *pb.Empty) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "hello"}, nil
}
