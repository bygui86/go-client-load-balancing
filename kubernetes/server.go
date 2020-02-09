package kubernetes

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bygui86/go-grpc/logger"

	"github.com/gorilla/mux"
)

const (
	KubeProbesStatusOk    = "OK"
	KubeProbesStatusError = "ERROR"

	KubeProbesCodeOk    = 200
	KubeProbesCodeError = 500
)

// KubeProbesServer - Kubernetes REST server
type KubeProbesServer struct {
	Config     *Config
	Router     *mux.Router
	HttpServer *http.Server
}

type KubeProbes struct {
	GrpcInterface GrpcInterface
}

// GrpcInterface -
type GrpcInterface interface {
	CheckState() (int, string, string)
}

var kubeProbes KubeProbes

// NewKubeProbesServer - Create new Kubernetes REST server
func NewKubeProbesServer(probes KubeProbes) (*KubeProbesServer, error) {
	logger.SugaredLogger.Info("Setup new Kubernetes probes server...")

	// create config
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}

	if !cfg.IsValid() {
		return nil, errors.New("Config not valid")
	}

	// create router
	router := newRouter()

	// create http server
	httpServer := newHttpServer(cfg.RestPort, router)

	// init KubeProbes object
	kubeProbes = probes

	// create KubeProbesServer object
	return &KubeProbesServer{
		Config:     cfg,
		Router:     router,
		HttpServer: httpServer,
	}, nil
}

// newRouter -
func newRouter() *mux.Router {
	logger.SugaredLogger.Debug("Setup new HTTP Router config...")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/live", livenessHandler)
	router.HandleFunc("/ready", readynessHandler)
	return router
}

// newHttpServer -
func newHttpServer(port int, router *mux.Router) *http.Server {
	logger.SugaredLogger.Debugf("Setup new HTTP server on port %d...", port)

	return &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(port),
		Handler: router,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
}

// Start - Start Kubernetes REST server
func (s *KubeProbesServer) Start() {
	logger.SugaredLogger.Info("Start Kubernetes probes server...")

	go func() {
		if err := s.HttpServer.ListenAndServe(); err != nil {
			logger.SugaredLogger.Error(err)
		}
	}()

	logger.SugaredLogger.Info("Kubernetes probes server listen on port ", s.Config.RestPort)
}

// Shutdown - Shutdown Kubernetes REST server
func (s *KubeProbesServer) Shutdown() {
	logger.SugaredLogger.Info("Shutdown Kubernetes probes server...")

	if s.HttpServer != nil {
		// create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.Config.ShutdownTimeout)*time.Second)
		defer cancel()
		// does not block if no connections, otherwise wait until the timeout deadline
		s.HttpServer.Shutdown(ctx)
	}
}
