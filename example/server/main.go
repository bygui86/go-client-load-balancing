package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"

	protoecho "google.golang.org/grpc/examples/features/proto/echo"
)

var (
	ipAddresses = []string{":50051", ":50052"}
)

type ecServer struct {
	protoecho.UnimplementedEchoServer // UnimplementedEchoServer can be embedded to have forward compatible implementations.
	ipAddress                         string
}

func main() {
	var wg sync.WaitGroup
	for _, ipAddress := range ipAddresses {
		wg.Add(1)
		go func(ipAddress string) {
			defer wg.Done()
			startServer(ipAddress)
		}(ipAddress)
	}
	wg.Wait()
}

func startServer(ipAddress string) {
	listener, listenErr := net.Listen("tcp", ipAddress)
	if listenErr != nil {
		log.Fatalf("Failed to listen: %v", listenErr)
	}

	grpcServer := grpc.NewServer()
	protoecho.RegisterEchoServer(grpcServer, &ecServer{ipAddress: ipAddress})
	log.Printf("gRPC serving on %s \n", ipAddress)
	serveErr := grpcServer.Serve(listener)
	if serveErr != nil {
		log.Fatalf("Failed to serve: %v", serveErr)
	}
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *protoecho.EchoRequest) (*protoecho.EchoResponse, error) {
	log.Printf("Received %s on %s \n", req.Message, s.ipAddress)
	return &protoecho.EchoResponse{
		Message: fmt.Sprintf("%s (from %s)", req.Message, s.ipAddress),
	}, nil
}
