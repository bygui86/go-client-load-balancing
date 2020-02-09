
# VARIABLES
# -


# CONFIG
.PHONY: help print-variables
.DEFAULT_GOAL := help


# ACTIONS

## code

build-protobuf :		## Compile protobuf
	protoc --proto_path=./proto/ --go_out=plugins=grpc:domain ./proto/*

build-server :		## Build server
	go build -o grpc-server ./hello-service

build-client :		## Build client
	go build -o grpc-client ./greeting-service

run-server :		## Run server
	@echo "Remember to run 'make build-protobuf' before starting the server"
	GO111MODULE=on go run ./hello-service/main.go

run-client :		## Run client
	@echo "Remember to run 'make build-protobuf' before starting the client"
	GO111MODULE=on go run ./greeting-service/main.go

## container

container-build-server :		## Build container image of the server
	docker build -t grpc/hello-service -f hello.Dockerfile .

container-build-client :		## Build container image of the client
	docker build -t grpc/greeting-service -f greeting.Dockerfile .

container-run-server :		## Run container of the server
	docker run -ti --rm --name hello-service -p 50051:50051 grpc/hello-service

container-run-client :		## Run container of the client
	docker run -ti --rm --name greeting-service grpc/greeting-service

## kubernetes

start-minikube :		## Start Minikube
	minikube start --cpus 4 --memory 8192 --disk-size=10g

start-kind :		## Start KinD
	kind create cluster --wait=60s

deploy-server :		## Deploy server on Kubernetes
	kubectl apply -k kube/hello-service

deploy-client :		## Deploy server on Kubernetes
	kubectl apply -k kube/greeting-service

client-logs :		## Show client logs
	kubectl logs -l app=greeting-service -f

## helpers

help :		## Help
	@echo ""
	@echo "*** \033[33mMakefile help\033[0m ***"
	@echo ""
	@echo "Targets list:"
	@grep -E '^[a-zA-Z_-]+ :.*?## .*$$' $(MAKEFILE_LIST) | sort -k 1,1 | awk 'BEGIN {FS = ":.*?## "}; {printf "\t\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""

print-variables :		## Print variables values
	@echo ""
	@echo "*** \033[33mMakefile variables\033[0m ***"
	@echo ""
	@echo "- - - makefile - - -"
	@echo "MAKE: $(MAKE)"
	@echo "MAKEFILES: $(MAKEFILES)"
	@echo "MAKEFILE_LIST: $(MAKEFILE_LIST)"
	@echo ""
