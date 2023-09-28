## Raft Consensus Algorithm with gRPC & Go

This project aims to demonstrate the working of the **Raft consensus** algorithm using gRPC and Go. The primary objective is
to implement a distributed consensus system where peers elect a leader and exchange messages to maintain consistency 
in a distributed network. Below are the key components and steps involved in the project:

## Project Setup,

To run the project, follow these steps:

Build the project:

```
$ make build
```

Run the peer

```
$ ./bin/peer
```

## Prerequisites

Before running the project, ensure you have the following prerequisites installed:

* Go v1.21
* Make cli
* protobuf-compiler

## Install Protobuf Tools

Install the necessary Protobuf tools with the following commands:

```
 go install http://google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
 go install http://google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

## To-Do List

The project is a work in progress, and the following tasks are planned for future development:

* Need to implement Raft election logic
* Admin cli to interact with peers
* Containerization
