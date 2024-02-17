BINARY_NAME=peer
LEADER_RPC_PORT=51500

run-leader:
	go run cmd/* --config node_configs/leader.yml

run-follower-1:
	go run cmd/* --config node_configs/follower-1.yml

run-follower-2:
	go run cmd/* --config node_configs/follower-2.yml

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protocol/cluster.proto

build:
	make gen-proto
	GOARCH=amd64 GOOS=darwin go build -pgo=auto -o bin/${BINARY_NAME}-darwin cmd/*
	GOARCH=amd64 GOOS=linux go build -pgo=auto -o bin/${BINARY_NAME}-linux cmd/*
	GOARCH=amd64 GOOS=windows go build -pgo=auto -o bin/${BINARY_NAME}-windows cmd/*
