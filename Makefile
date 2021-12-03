.PHONY: build
build : buildServer buildClient

.PHONY: server buildServer
server buildServer : bin/server

bin/server: protobuf/singularity6.pb.go protobuf/singularity6_grpc.pb.go server/main.go
	go build -o bin/server server/main.go

.PHONY: client buildClient
client buildClient : bin/client

bin/client: protobuf/singularity6.pb.go protobuf/singularity6_grpc.pb.go client/main.go
	go build -o bin/client client/main.go

.PHONY: proto buildProto
proto buildProto : protobuf/singularity6.pb.go protobuf/singularity6_grpc.pb.go	

protobuf/singularity6.pb.go protobuf/singularity6_grpc.pb.go: protobuf/singularity6.proto
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protobuf/singularity6.proto

.PHONY: all
all: build

.PHONY: clean
clean:
	rm bin/*
	rm protobuf/*.pb.go
	go clean
