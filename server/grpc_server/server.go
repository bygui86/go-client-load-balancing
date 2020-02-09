package grpc_server

import (
	"context"
	"github.com/bygui86/go-grpc-client-lb/domain"
	"github.com/bygui86/go-grpc-client-lb/logger"
)

// Server - Used to implement domain.EchoService
type Server struct{}

// Echo - Implement service EchoService.Echo
func (s *Server) Echo(ctx context.Context, in *domain.EchoRequest) (*domain.EchoResponse, error) {
	logger.SugaredLogger.Infof("Echo message: %s", in.Message)
	return &domain.EchoResponse{
		Message: in.Message,
	}, nil
}
