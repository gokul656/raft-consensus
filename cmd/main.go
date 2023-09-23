package main

import (
	"fmt"
	"log"
	"net"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	conf := config.GetEnv()
	rpcPort := fmt.Sprintf(":%s", conf.RPCPort)
	apiPort := fmt.Sprintf(":%s", conf.APIPort)

	go setupRest(rpcPort)
	go setupRPC(apiPort)

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
		panic(common.INVALID_RPC_PORT)
	}

	if err := server.Serve(listen); err != nil {
		panic(err)
	}

	reflection.Register(server)
}

func setupRest(port string) {
	// TODO : Implement restful APIs
}
