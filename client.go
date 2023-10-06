package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gokul656/raft-consensus/protocol"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:51505", opts)
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	client := protocol.NewClusterClient(cc)
	resp, err := client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive response => %s ", resp.Message)
}
