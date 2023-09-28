run-peer:
	go run cmd/* --api-port 3001

gen-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protocol/cluster.proto

build:
	go build -o bin/peer cmd/*
