
### Protobuf compiling stage
FROM golang:1.13-buster AS protobuilder

ARG PROTOC_VERSION=3.9.1

# Install dependencies
RUN apt-get update \
	&& apt-get install -y unzip wget

# Install protobuf environment
WORKDIR /opt/protoc
RUN wget -q https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/protoc-$PROTOC_VERSION-linux-x86_64.zip -O protoc.zip \
	&& unzip -o protoc.zip -d /usr/local bin/protoc \
	&& unzip -o protoc.zip -d /usr/local/ include/* \
	&& rm protoc.zip
RUN go get -u google.golang.org/grpc \
	&& go get -u github.com/golang/protobuf/protoc-gen-go \
	&& go get -u github.com/square/goprotowrap/cmd/protowrap

WORKDIR /app

# Copy protobuf
COPY proto proto

# Compile protobuf
RUN mkdir -p domain/
RUN protoc --proto_path=./proto/ --go_out=plugins=grpc:domain ./proto/*

### ---

### Code compiling stage
FROM protobuilder AS gobuilder

ENV GO111MODULE=on

WORKDIR /github.com/bygui86/go-grpc

# Prepare go-modules environment
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy compiled protobuf
COPY --from=protobuilder /app/domain domain

# Copy application code
COPY client client
COPY kubernetes kubernetes
COPY logger logger
COPY utils utils

# Compile application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/app ./client/main.go

### --- ALPINE

### Final image
FROM alpine

# Install additionals
RUN apk add --no-cache bash

# Copy application executable
COPY --from=gobuilder /bin/app /bin/app

# gRPC
EXPOSE 50051
# KubeProbes
EXPOSE 9091

# Run application
ENTRYPOINT ["/bin/app"]

### --- SCRATCH

#### Final image
#FROM scratch
#
## Copy application executable
#COPY --from=gobuilder /bin/app /bin/app
#
### gRPC
#EXPOSE 50051
## KubeProbes
#EXPOSE 9091
#
#ENTRYPOINT ["/bin/app"]
