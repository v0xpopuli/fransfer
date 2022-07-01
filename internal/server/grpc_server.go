package server

import (
	"fransfer/internal/auth"
	"fransfer/internal/filetransfer"
	"fransfer/internal/generated"
	"net"

	"google.golang.org/grpc"
)

type (
	Server interface {
		Register(filetransfer.Service)
		Run() error
	}

	server struct {
		srv  *grpc.Server
		addr string
	}
)

func New(addr string, jwt auth.JWT) Server {
	return server{
		addr: addr,
		srv:  grpc.NewServer(WithStreamServerAuthorizationInterceptor(jwt)),
	}
}

func (s server) Register(service filetransfer.Service) {
	generated.RegisterFileTransferServer(s.srv, service)
}

func (s server) Run() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.srv.Serve(l)
}
