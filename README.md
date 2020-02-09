
# go-grpc-client-lb
Example of gRPC client-side load balancing in Golang

## Services

- [gRPC server](server)
- [gRPC client](client)

---

## Build

1. Server
	```shell
    make build-server
    # or
	go build -o grpc-server ./server
	```

2. Client
	```shell
    make build-client
    # or
	go build -o grpc-client ./client
	```

---

## Run

1. Run server
	```shell
    make run-server
    # or
	./grpc-server
	```

2. In another shell, run client
	```shell
    make run-client
    # or
	./grpc-client
	```

### From source

1. Run server
	```shell
    make run-server-src
    # or
    protoc --proto_path=./proto/ --go_out=plugins=grpc:domain ./proto/*
	GO111MODULE=on go run ./server/main.go
	```

2. In another shell, run client
	```shell
    make run-client-src
    # or
    protoc --proto_path=./proto/ --go_out=plugins=grpc:domain ./proto/*
	GO111MODULE=on go run ./client/main.go
	```

---

## Docker

### Build

1. Server
	```shell
    make container-build-server
    # or
	docker build -t grpc/grpc-server -f server.Dockerfile .
	```

2. Client
	```shell
    make container-build-client
    # or
	docker build -t grpc/grpc-client -f client.Dockerfile .
	```

### Run

1. Run server
	```shell
    make container-run-server
    # or
	docker run -ti --rm --name server -p 50051:50051 grpc/grpc-server
	```

2. In another shell, run client
	```shell
    make container-run-client
    # or
	docker run -ti --rm --name client grpc/grpc-client
	```

---

## Kubernetes

### Import container images

#### Minikube

1. Enable Minikube internal container registry
    ```shell
    eval $(minikube docker-env)
    ```

2. Build container images normally as it would be locally

#### KinD

1. Build container images normally as it would be locally

2. Import in KinD
    ```shell
    make load-container-kind
    # or
    kind load docker-image grpc/grpc-server
    kind load docker-image grpc/grpc-client
    ```

### Deploy

1. Start Kubernetes locally

    - Minikube
        ```shell
        make start-minikube
        # or
        minikube start --cpus 4 --memory 8192 --disk-size=10g
        ```

    - Kind
        ```shell
        make start-kind
        # or
        kind create cluster --wait=60s
        ```

2. Deploy server
    ```shell
    make deploy-server
    # or
    kubectl apply -k kube/server
    ```

3. Deploy client
    ```shell
    make deploy-client
    # or
    kubectl apply -k kube/client
    ```

4. Take a look of logs
    - Server
        ```shell
        make server-logs
        # or
        kubectl logs -l app=grpc-server -f
        ```
    - (In another shell) Client
        ```shell
        make client-logs
        # or
        kubectl logs -l app=grpc-client -f
        ```

### Cleanup

#### Minikube
```shell
make stop-minikube
# or
minikube stop
minikube delete
```

#### KinD
```shell
make stop-kind
# or
kind delete cluster
```

---

## Example

1. Run server
    ```shell
    go run example/server/main.go
    ```

2. In another shell, run client
    ```shell
    go run example/client/main.go
    ```

---

## Links
- https://grpc.io/blog/loadbalancing/
- https://github.com/grpc/grpc-go/tree/master/examples/features/load_balancing
- https://github.com/grpc/grpc/blob/master/doc/service_config.md
- https://medium.com/@ammar.daniel/grpc-client-side-load-balancing-in-go-cd2378b69242
- https://www.marwan.io/blog/grpc-dns-load-balancing
