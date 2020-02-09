package grpc_server

import (
	"fmt"
	"net"

	"github.com/bygui86/go-grpc-client-lb/kubernetes"
)

type GrpcServerService struct {
	Listener net.Listener
	Network  string
	Address  string
}

func (s *GrpcServerService) CheckState() (int, string, string) {
	address := s.Listener.Addr()

	if address.Network() != s.Network {
		return kubernetes.KubeProbesCodeError,
			kubernetes.KubeProbesStatusError,
			fmt.Sprintf("gRPC listener network mismatch: expected[%s] actual[%s]", s.Network, address.Network())
	}

	if address.String() != s.Address &&
		address.String() != "[::]:50051" {
		return kubernetes.KubeProbesCodeError,
			kubernetes.KubeProbesStatusError,
			fmt.Sprintf("gRPC listener address mismatch: expected[%s] actual[%s]", s.Address, address.String())
	}

	return kubernetes.KubeProbesCodeOk,
		kubernetes.KubeProbesStatusOk,
		fmt.Sprintf("gRPC server listening on %s", address.String())
}
