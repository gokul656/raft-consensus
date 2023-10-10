package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
	"github.com/gokul656/raft-consensus/internal"
	"github.com/gokul656/raft-consensus/protocol"
	"github.com/gokul656/raft-consensus/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	conf := config.GetEnv()
	rpcPort := fmt.Sprintf(":%s", conf.RPCPort)
	apiPort := fmt.Sprintf(":%s", conf.APIPort)

	go setupRPC(rpcPort)
	go setupRest(apiPort)

	log.Println("[ rpc ] listening at         ", rpcPort)
	log.Println("[ api ] listening at         ", apiPort)
	log.Println("[ log ] files can be found at", conf.LogDir)

	select {}
}

func setupRPC(port string) {
	defer common.HandlePanic("setupRPC")

	server := grpc.NewServer()
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln("bind address already in use")
	}

	internal.StartupRaft()
	protocol.RegisterClusterServer(server, rpc.NewgRPCServer())
	reflection.Register(server)

	if err = server.Serve(listen); err != nil {
		panic(err)
	}
}

func setupRest(port string) {
	// TODO : Implement restful APIs
}
