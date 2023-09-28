package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
	"github.com/gokul656/raft-consensus/protocol"
	"github.com/gokul656/raft-consensus/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	conf := config.GetEnv()
	rpcPort := fmt.Sprintf(":%s", conf.RPCPort)
	apiPort := fmt.Sprintf(":%s", conf.APIPort)

	go setupRest(apiPort)
	go setupRPC(rpcPort)

	log.Println("[ rpc ] listening at         ", rpcPort)
	log.Println("[ api ] listening at         ", apiPort)
	log.Println("[ log ] files can be found at", conf.TmpDir)

	select {}
}

func setupRPC(port string) {
	defer common.HandlePanic()

	server := grpc.NewServer()
	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(common.InvalidRPCPort)
	}

	protocol.RegisterClusterServer(server, rpc.NewgRPCServer())
	reflection.Register(server)

	if err = server.Serve(listen); err != nil {
		panic(err)
	}
}

func setupRest(port string) {
	// TODO : Implement restful APIs
}
