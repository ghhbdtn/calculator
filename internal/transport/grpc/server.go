package grpc

import (
	"log"
	"net"

	"calculator/internal/app/service"
	"calculator/pkg/calculatorpb"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	calculator *service.Calculator
}

func NewServer(calculator *service.Calculator) *Server {
	s := &Server{
		grpcServer: grpc.NewServer(),
		calculator: calculator,
	}

	handler := NewHandler(calculator)
	calculatorpb.RegisterCalculatorServiceServer(s.grpcServer, handler)

	return s
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Printf("gRPC server starting on %s", addr)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
