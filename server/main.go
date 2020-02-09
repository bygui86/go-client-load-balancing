package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bygui86/go-grpc-client-lb/domain"
	"github.com/bygui86/go-grpc-client-lb/kubernetes"
	"github.com/bygui86/go-grpc-client-lb/logger"
	"github.com/bygui86/go-grpc-client-lb/server/grpc_server"
	"github.com/bygui86/go-grpc-client-lb/utils"

	"google.golang.org/grpc"
)

const (
	serverAddressEnvVar  = "GOGRPC_GRPC_SERVER_ADDRESS"
	kubeProbesNameEnvVar = "GOGRPC_KUBE_PROBES_START"

	serverAddressEnvVarDefault  = "0.0.0.0:50051"
	kubeProbesNameEnvVarDefault = false

	grpcListenerNetwork = "tcp"
)

func main() {
	address := utils.GetString(serverAddressEnvVar, serverAddressEnvVarDefault)
	kubeProbes := utils.GetBool(kubeProbesNameEnvVar, kubeProbesNameEnvVarDefault)

	listener := createListener(grpcListenerNetwork, address)
	defer listener.Close()
	logger.SugaredLogger.Infof("TCP listener ready on %s", address)

	go startGrpcServer(listener)
	logger.SugaredLogger.Infof("gRPC server ready")

	if kubeProbes {
		kubeServer := startKubernetes(listener, grpcListenerNetwork, address)
		defer kubeServer.Shutdown()
	}

	logger.SugaredLogger.Info("gRPC server started!")
	startSysCallChannel()
}

// createListener -
func createListener(network, address string) net.Listener {
	listener, err := net.Listen(network, address)
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to listen: %v", err.Error())
		os.Exit(3)
	}

	return listener
}

// startGrpcServer -
func startGrpcServer(listener net.Listener) {
	grpcServer := grpc.NewServer()
	helloSvcServer := grpc_server.Server{}
	domain.RegisterEchoServiceServer(grpcServer, &helloSvcServer)
	err := grpcServer.Serve(listener)
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to serve: %v", err)
		os.Exit(4)
	}
}

// startKubernetes -
func startKubernetes(listener net.Listener, network, address string) *kubernetes.KubeProbesServer {
	kubeProbes := kubernetes.KubeProbes{
		GrpcInterface: &grpc_server.GrpcServerService{
			Listener: listener,
			Network:  network,
			Address:  address,
		},
	}
	server, err := kubernetes.NewKubeProbesServer(kubeProbes)
	if err != nil {
		logger.SugaredLogger.Errorf("Kubernetes probes server creation failed: %s", err.Error())
		os.Exit(2)
	}
	logger.SugaredLogger.Debug("Kubernetes probes server successfully created")

	server.Start()
	logger.SugaredLogger.Debug("Kubernetes probes successfully started")

	return server
}

// startSysCallChannel -
func startSysCallChannel() {
	syscallCh := make(chan os.Signal)
	signal.Notify(syscallCh, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-syscallCh
	logger.SugaredLogger.Info("Termination signal received!")
}
