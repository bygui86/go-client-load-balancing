package kubernetes

import (
	"encoding/json"
	"net/http"

	"github.com/bygui86/go-grpc/logger"
)

// livenessHandler -
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	logger.SugaredLogger.Debug("Liveness probe invoked...")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		Liveness{
			Code:   KubeProbesCodeOk,
			Status: KubeProbesStatusOk,
		},
	)
}

// readynessHandler -
func readynessHandler(w http.ResponseWriter, r *http.Request) {
	logger.SugaredLogger.Debug("Readyness probe invoked...")

	generalCode := KubeProbesCodeOk
	generalStatus := KubeProbesStatusOk

	grpcReady := checkGrpc()
	if grpcReady.Code == KubeProbesCodeError {
		generalCode = KubeProbesCodeError
		generalStatus = KubeProbesStatusError
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		Readyness{
			Code:   generalCode,
			Status: generalStatus,
			Grpc:   grpcReady,
		},
	)
}

// checkGrpc -
func checkGrpc() GrpcReadyness {
	code := KubeProbesCodeOk
	status := KubeProbesStatusOk
	msg := ""

	if kubeProbes.GrpcInterface != nil {
		code, status, msg = kubeProbes.GrpcInterface.CheckState()
	}

	return GrpcReadyness{
		Code:    code,
		Status:  status,
		Message: msg,
	}
}
