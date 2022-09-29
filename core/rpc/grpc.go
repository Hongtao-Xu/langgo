package rpc

import (
	"context"
	ctls "crypto/tls"
	"crypto/x509"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"langgo/core/log"
	"net"
	"runtime/debug"
	"time"
)

type Server struct {
	opt        []grpc.ServerOption
	middleware []grpc.UnaryServerInterceptor
	server     *grpc.Server
	tls        *Tls
}

type Tls struct {
	Crt   string
	Key   string
	CACrt string
}

//Use 拦截器
func (s *Server) Use(middleware ...grpc.UnaryServerInterceptor) {
	s.middleware = append(s.middleware, middleware...)
}

//NewClient 构造rpc客户端
func NewClient(tls *Tls, addr string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	if tls != nil {
		certificate, err := ctls.LoadX509KeyPair(tls.Crt, tls.Key)
		if err != nil {
			return nil, err
		}
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(tls.CACrt)
		if err != nil {
			return nil, err
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, errors.New("AppendCertsFromPEM is false")
		}
		creds := credentials.NewTLS(&ctls.Config{
			Certificates:       []ctls.Certificate{certificate},
			RootCAs:            certPool,
			InsecureSkipVerify: true,
		})

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	return grpc.Dial(addr, opts...)
}

//NewServer 构造rpc服务端
func NewServer(tls *Tls, opt ...grpc.ServerOption) *Server {
	s := &Server{
		opt: opt,
		tls: tls,
	}
	return s
}

// Server 添加Server服务,tls
func (s *Server) Server() (server *grpc.Server, err error) {
	if s.server == nil {
		s.opt = append(s.opt, grpc.UnaryInterceptor(ChainUnaryServer(s.middleware...)))
		if s.tls != nil {
			//生成证书
			certificate, err := ctls.LoadX509KeyPair(s.tls.Crt, s.tls.Key)
			if err != nil {
				return nil, err
			}
			//证书池
			certPool := x509.NewCertPool()
			ca, err := ioutil.ReadFile(s.tls.CACrt)
			if err != nil {
				return nil, err
			}
			if ok := certPool.AppendCertsFromPEM(ca); !ok { //证书池加入ca
				panic("AppendCertsFromPEM failed")
			}
			creds := credentials.NewTLS(&ctls.Config{
				Certificates: []ctls.Certificate{certificate},
				ClientAuth:   ctls.RequireAndVerifyClientCert,
				ClientCAs:    certPool})
			s.opt = append(s.opt, grpc.Creds(creds))
		}
		s.server = grpc.NewServer(s.opt...)
	}
	return s.server, nil
}

//Run grpc的调用
func (s *Server) Run(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Logger("grpc", "run").Info().Str("addr", addr).Msg("server ready")
	if err := s.server.Serve(lis); err != nil {
		log.Logger("grpc", "run").Error().Err(err).Send()
		return err
	}
	return nil
}

//ChainUnaryServer 链式一元拦截器（多拦截器）
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)

	//Dummy 拦截器，避免返回nil
	if n == 0 {
		return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			return handler(ctx, req)
		}
	}
	//直接返回单个拦截器
	if n == 1 {
		return interceptors[0]
	}
	//返回拦截器接口的函数
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		currHandler := handler

		for i := n - 1; i > 0; i-- {
			innerHandler, i := currHandler, i
			currHandler = func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				//链式传递
				return interceptors[i](currentCtx, currentReq, info, innerHandler)
			}
		}
		//最后一个拦截器
		return interceptors[0](ctx, req, info, currHandler)
	}
}

// LogUnaryServerInterceptor 日志拦截器
func LogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		st := time.Now()
		defer func() {
			//如果有异常的日志信息
			if recoverError := recover(); recoverError != nil {
				log.Logger("grpc", "log").
					Error().Interface("recover", recoverError).
					Bytes("stack", debug.Stack()).
					Interface("req", req).
					Interface("resp", resp).
					Err(err).TimeDiff("runtime", time.Now(), st).
					Send()
			} else { //正常调用的日志信息
				log.Logger("grpc", "log").
					Info().Interface("req", req).
					Interface("resp", resp).
					Err(err).TimeDiff("runtime", time.Now(), st).
					Send()
			}
		}()
		//链式传递
		resp, err = handler(ctx, req)
		return resp, err
	}
}
