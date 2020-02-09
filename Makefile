
# VARIABLES
# -


# CONFIG
.PHONY: help print-variables
.DEFAULT_GOAL := help


# ACTIONS

## code

build-proto :		## Compile protobuf
	protoc --proto_path=./proto/ --go_out=plugins=grpc:domain ./proto/*

build-server : build-proto		## Build server
	GO111MODULE=on go build -o grpc-server ./server

build-client : build-proto		## Build client
	GO111MODULE=on go build -o grpc-client ./client

run-server-src : build-proto		## Run server
	GO111MODULE=on go run ./server/main.go

run-client-src : build-proto		## Run client
	GO111MODULE=on go run ./client/main.go

run-server : build-server		## Run server
	./grpc-server

run-client : build-client		## Run client
	./grpc-client

## container

container-build-server :		## Build container image of the server
	docker build -t grpc/grpc-server -f server.Dockerfile .

container-build-client :		## Build container image of the client
	docker build -t grpc/grpc-client -f client.Dockerfile .

container-run-server :		## Run container of the server
	docker run -ti --rm --name server -p 50051:50051 grpc/server

container-run-client :		## Run container of the client
	docker run -ti --rm --name client grpc/client

## kubernetes

start-minikube :		## Start Minikube
	minikube start --cpus 4 --memory 8192 --disk-size=10g

start-kind :		## Start KinD
	kind create cluster --wait=60s

stop-minikube :		## Stop Minikube
	minikube stop
	minikube delete

stop-kind :		## Stop KinD
	kind delete cluster

load-container-kind :		## Load container images in KinD
	kind load docker-image grpc/grpc-server
	kind load docker-image grpc/grpc-client

deploy-server :		## Deploy server on Kubernetes
	kubectl apply -k kube/server

deploy-client :		## Deploy server on Kubernetes
	kubectl apply -k kube/client

delete-server :		## Delete server from Kubernetes
	kubectl delete -k kube/server

delete-client :		## Delete server from Kubernetes
	kubectl delete -k kube/client

server-logs :		## Show server logs
	kubectl logs -l app=grpc-server -f

client-logs :		## Show client logs
	kubectl logs -l app=grpc-client -f

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
