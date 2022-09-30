package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/Hongtao-Xu/langgo"
	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/core/rpc"
	cs "github.com/Hongtao-Xu/langgo/examples/grpc/etcd/server/components/server"
	"github.com/Hongtao-Xu/langgo/examples/grpc/etcd/server/pb"
	"github.com/Hongtao-Xu/langgo/examples/grpc/etcd/server/service/server"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 8001, "port")
	flag.Parse()
	addr := fmt.Sprintf("localhost:%d", port)
	//1.加载Instance
	langgo.Run(&cs.Instance{})
	//2.关闭时，etcd移除
	core.SignalHandle(&core.SignalHandler{
		Sig: syscall.SIGINT,
		F: func() {
			rpc.EtcdUnRegister(cs.GetInstance().ServiceName, addr)
			os.Exit(int(syscall.SIGINT))
		},
	})
	//deferRun
	defer func() {
		core.DeferRun()
	}()
	//3.注册etcd
	rpc.EtcdRegister(cs.GetInstance().EtcdHost, cs.GetInstance().ServiceName, addr, 50)
	cg := rpc.NewServer(nil)
	cg.Use(rpc.LogUnaryServerInterceptor())
	//4.生成gs
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
