package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Hongtao-Xu/langgo/core/rpc"
	"github.com/Hongtao-Xu/langgo/examples/grpc/etcd/client/pb"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

const etcdHost = "http://localhost:2379"
const serviceName = "langgo/server"

func main() {
	//1.etcd相关配置
	etcdClient, err := clientv3.NewFromURL(etcdHost)
	if err != nil {
		panic(err)
	}
	etcdResolver, err := resolver.NewBuilder(etcdClient)
	//2.grpc相关配置
	conn, err := rpc.NewClient(nil, fmt.Sprintf("etcd:///%s", serviceName),
		grpc.WithResolvers(etcdResolver),                                                             //配置etcd
		grpc.WithTransportCredentials(insecure.NewCredentials()),                                     //配置凭证
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name))) //负载均衡
	if err != nil {
		panic(err)
	}
	serverClient := pb.NewServerClient(conn)
	//3.循环访问Hello接口
	for {
		helloRespone, err := serverClient.Hello(context.Background(), &pb.Empty{})
		if err != nil {
			fmt.Printf("err: %v", err)
			return
		}
		log.Println(helloRespone, err)
		time.Sleep(500 * time.Millisecond)
	}
}
