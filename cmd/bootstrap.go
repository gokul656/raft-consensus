package main

import (
	"context"
	"log"

	"github.com/gokul656/raft-consensus/peer"
)

func init() {
	log.Println("Setting up cluster...")
	cluster := &peer.RaftHub{
		Ctx: context.Background(),
		Self: &peer.Peer{
			Address: "localhost:3001",
			Name:    "leader",
			State:   peer.LEADER,
		},
		Peers: peer.NewMap[string, *peer.Peer](),
	}

	cluster.Leader = cluster.Self
	cluster.Peers.Put("client", &peer.Peer{
		Name:    "follower",
		Address: "localhost:3001",
		State:   peer.FOLLOWER,
	})

	// cluster.InitiateElection()
}
