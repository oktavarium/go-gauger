package grpcserver

import (
	"fmt"
	"net"

	pbapi "github.com/oktavarium/go-gauger/api"
	"github.com/oktavarium/go-gauger/internal/server/internal/gaugeserver/internal/storage"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	pbapi.UnimplementedGaugerServer
	addr    string
	storage storage.Storage
	server  *grpc.Server
}

func NewGrpcServer(addr string, s storage.Storage) *GrpcServer {
	return &GrpcServer{
		addr:    addr,
		storage: s,
		server:  grpc.NewServer(),
	}
}

func (s *GrpcServer) ListenAndServe() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("error on listening socket for grpc: %w", err)
	}

	// регистрируем сервис
	pbapi.RegisterGaugerServer(s.server, s)

	if err := s.server.Serve(listen); err != nil {
		return fmt.Errorf("error on serving grpc: %w", err)
	}

	return nil
}

func (s *GrpcServer) Shutdown() {
	s.server.GracefulStop()
}
