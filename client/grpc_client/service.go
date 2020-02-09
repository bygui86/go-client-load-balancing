package grpc_client

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

	"github.com/bygui86/go-grpc-client-lb/kubernetes"
)

type GrpcClientService struct {
	GrpcClientConn *grpc.ClientConn
}

func (s *GrpcClientService) CheckState() (int, string, string) {
	target := s.GrpcClientConn.Target()
	state := s.GrpcClientConn.GetState()

	if target == "" {
		return kubernetes.KubeProbesCodeError,
			kubernetes.KubeProbesStatusError,
			"gRPC target not set"
	}

	if state != connectivity.Ready {
		return kubernetes.KubeProbesCodeError,
			kubernetes.KubeProbesStatusError,
			fmt.Sprintf("gRPC connection state not ready: %s", state)
	}

	return kubernetes.KubeProbesCodeOk,
		kubernetes.KubeProbesStatusOk,
		fmt.Sprintf("gRPC connected to %s", target)
}
