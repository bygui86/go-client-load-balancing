package kubernetes

import (
	"github.com/bygui86/go-grpc/logger"
	"github.com/bygui86/go-grpc/utils"
)

const (
	restPortEnvVar                  = "GOGRPC_KUBE_PROBES_PORT"
	restServerShutdownTimeoutEnvVar = "GOGRPC_KUBE_SERVER_SHUTDOWN_TIMEOUT"

	restPortDefault           = 9091
	restServerShutdownTimeout = 10
)

// Config -
type Config struct {
	RestPort        int
	ShutdownTimeout int
}

// IsValid -
func (c *Config) IsValid() bool {
	return c.RestPort > 0
}

// newConfig -
func newConfig() (*Config, error) {
	logger.SugaredLogger.Debug("Setup new Kubernetes config...")

	return &Config{
		RestPort:        utils.GetInt(restPortEnvVar, restPortDefault),
		ShutdownTimeout: utils.GetInt(restServerShutdownTimeoutEnvVar, restServerShutdownTimeout),
	}, nil
}
