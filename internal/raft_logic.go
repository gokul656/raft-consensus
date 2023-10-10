package internal

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
	"github.com/gokul656/raft-consensus/peer"
	"github.com/gokul656/raft-consensus/protocol"
)

var Cluster *peer.RaftHub
var RPCServer protocol.ClusterServer
var once sync.Once

func StartupRaft() {
	log.Println("Setting up cluster...")
	once.Do(
		func() {
			defer common.HandlePanic("raft_state")

			env := config.GetEnv()
			Cluster = peer.NewRaft(&peer.Peer{
				Address: fmt.Sprintf("localhost:%s", env.RPCPort),
				Name:    env.InstanceID,
				State:   protocol.PeerState_FOLLOWER.Enum(),
			})

			Cluster.AddPeer(env.InstanceID, fmt.Sprintf("localhost:%s", env.RPCPort))
			Cluster.Self = Cluster.GetPeer(env.InstanceID)

			// if there are no leaders, the current peer elects itself as leader & initiates election
			if env.Leader == "" {
				runAsLeader()
			} else {
				runAsFollower()
			}
		},
	)
}

func runAsLeader() {
	env := config.GetEnv()
	Cluster.ChangeLeader(env.InstanceID)

	go Cluster.CheckFollowersHealth()
}

func runAsFollower() {
	env := config.GetEnv()

	Cluster.AddPeer(env.LeaderID, env.Leader)
	Cluster.ChangeLeader(env.LeaderID)

	err := Cluster.Register(context.Background(), &protocol.AddPeerRequest{
		Peer: &protocol.Peer{
			Name:    Cluster.Self.Name,
			Address: Cluster.Self.Address,
		},
	})

	if err != nil {
		log.Fatalln("[CRITICAL] Unable to register as Follower", err)
	}

	Cluster.Synchronize()
	go Cluster.CheckLeaderHealth()
}
