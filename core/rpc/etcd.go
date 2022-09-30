package rpc

import (
	"context"
	"fmt"

	"github.com/Hongtao-Xu/langgo/core"
	"github.com/Hongtao-Xu/langgo/core/log"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
)

var etcdClient *clientv3.Client

// EtcdRegister etcd注册
func EtcdRegister(etcdHost, serviceName, addr string, ttl int64) error {
	//生成etcdClient
	etcdClient, err := clientv3.NewFromURL(etcdHost)
	if err != nil {
		return err
	}
	//生成租约
	lease, _ := etcdClient.Grant(context.TODO(), ttl)
	//生成em
	em, err := endpoints.NewManager(etcdClient, serviceName)
	if err != nil {
		return err
	}
	//em添加节点
	err = em.AddEndpoint(context.TODO(),
		fmt.Sprintf("%s/%s", serviceName, addr),
		endpoints.Endpoint{Addr: addr},
		clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	//调用链，取消注册
	core.DeferAdd(func() {
		EtcdUnRegister(serviceName, addr)
	})
	//etcd保持连接
	alive, err := etcdClient.KeepAlive(context.TODO(), lease.ID)
	if err != nil {
		return err
	}

	go func() { //监听etcd连接存活状态
		for {
			<-alive
			log.Logger("grpc", "etcd").Debug().Msg("Keep Alive")
		}
	}()
	return nil
}

// EtcdUnRegister etcd删除
func EtcdUnRegister(serviceName, addr string) error {
	log.Logger("grpc", "etcd").Debug().Str("addr", addr).Msg("unregister")
	if etcdClient != nil {
		//生成em
		em, err := endpoints.NewManager(etcdClient, serviceName)
		if err != nil {
			return err
		}
		//删除节点
		err = em.DeleteEndpoint(context.TODO(), fmt.Sprintf("%s/%s", serviceName, addr))
		if err != nil {
			return err
		}
		return err
	}
	return nil
}
