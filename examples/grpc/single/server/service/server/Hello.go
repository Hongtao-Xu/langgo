package server

import (
	"context"
	"langgo/examples/grpc/single/server/pb"
)

func (s Server) Hello(ctx context.Context, empty *pb.Empty) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "hello"}, nil
}
