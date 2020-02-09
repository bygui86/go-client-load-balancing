package kubernetes

type Liveness struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type Readyness struct {
	Code   int           `json:"code"`
	Status string        `json:"status"`
	Grpc   GrpcReadyness `json:"grpc"`
}

type GrpcReadyness struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
