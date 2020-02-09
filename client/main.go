package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bygui86/go-grpc-client-lb/client/grpc_client"
	"github.com/bygui86/go-grpc-client-lb/domain"
	"github.com/bygui86/go-grpc-client-lb/kubernetes"
	"github.com/bygui86/go-grpc-client-lb/logger"
	"github.com/bygui86/go-grpc-client-lb/utils"

	"google.golang.org/grpc"
)

const (
	serverAddressEnvVar  = "GOGRPC_SERVER_ADDRESS"
	messageEnvVar        = "GOGRPC_MESSAGE"
	kubeProbesNameEnvVar = "GOGRPC_KUBE_PROBES_START"

	serverAddressEnvVarDefault  = "0.0.0.0:50051"
	messageEnvVarDefault        = "Default message"
	kubeProbesNameEnvVarDefault = false

	// Available values: passthrough | dns
	grpcResolverScheme = "dns"
	grpcDefaultServiceConfig = `{
	"loadBalancingPolicy": "%s"
}
`
)

func main() {
	serverAddress := utils.GetString(serverAddressEnvVar, serverAddressEnvVarDefault)
	message := utils.GetString(messageEnvVar, messageEnvVarDefault)
	kubeProbes := utils.GetBool(kubeProbesNameEnvVar, kubeProbesNameEnvVarDefault)

	grpcConn := createGrpcConnection(serverAddress)
	defer grpcConn.Close()
	logger.SugaredLogger.Infof("gRPC connection ready to %s", serverAddress)

	go startMessageSender(grpcConn, message)

	if kubeProbes {
		kubeServer := startKubernetes(grpcConn)
		defer kubeServer.Shutdown()
	}

	logger.SugaredLogger.Info("gRPC client started!")
	startSysCallChannel()
}

// createGrpcConnection -
func createGrpcConnection(host string) *grpc.ClientConn {
	resolver.SetDefaultScheme(grpcResolverScheme)

	connection, err := grpc.Dial(
		host,
		//grpc.WithBalancerName(roundrobin.Name), // [Deprecated] This sets the initial balancing policy
		grpc.WithDefaultServiceConfig(
			fmt.Sprintf(grpcDefaultServiceConfig, roundrobin.Name),
		),
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)

	if err != nil {
		logger.SugaredLogger.Errorf("Connection to gRPC server %s failed: %v", host, err.Error())
		os.Exit(3)
	}

	logger.SugaredLogger.Info("State: ", connection.GetState())
	logger.SugaredLogger.Info("Target: ", connection.Target())

	return connection
}

// startMessageSender -
func startMessageSender(connection *grpc.ClientConn, message string) {
	timeout := 2 * time.Second
	client := domain.NewEchoServiceClient(connection)
	logger.SugaredLogger.Info("Starting message sender...")
	for {
		go sendMessage(client, timeout, message)
		time.Sleep(3 * time.Second)
	}
}

// sendMessage -
func sendMessage(client domain.EchoServiceClient, timeout time.Duration, message string) {
	// WARNING: the connection context is one-shot, it must be refreshed before every request
	connectionCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := client.Echo(connectionCtx, &domain.EchoRequest{Message: message})
	if err != nil {
		logger.SugaredLogger.Errorf("Could not send message %s: %v", message, err.Error())
		return
	}
	logger.SugaredLogger.Info(response.Message)
}

// startKubernetes -
func startKubernetes(grpcConn *grpc.ClientConn) *kubernetes.KubeProbesServer {
	kubeProbes := kubernetes.KubeProbes{
		GrpcInterface: &grpc_client.GrpcClientService{
			GrpcClientConn: grpcConn,
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
