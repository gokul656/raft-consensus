
BINARY_NAME=peer
LEADER_RPC_PORT=51500

run-leader:
	go run cmd/* --api-port 3001 --grpc-port ${LEADER_RPC_PORT}  --name peer-0

run-follower-1:
	go run cmd/* --api-port 3001 --grpc-port 51501 --leader localhost:${LEADER_RPC_PORT} --name peer-1

run-follower-2:
	go run cmd/* --api-port 3001 --grpc-port 51502 --leader localhost:${LEADER_RPC_PORT} --name peer-2

run-leader-and-follower:
	make run-leader &1
	make run-follower

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protocol/cluster.proto

build:
	make gen-proto
	GOARCH=amd64 GOOS=darwin go build -pgo=auto -o bin/${BINARY_NAME}-darwin cmd/*
	GOARCH=amd64 GOOS=linux go build -pgo=auto -o bin/${BINARY_NAME}-linux cmd/*
	GOARCH=amd64 GOOS=windows go build -pgo=auto -o bin/${BINARY_NAME}-windows cmd/*
