package peer

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Vote struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Term    uint   `json:"term"`
}

type RaftHub struct {
	Self     *Peer
	Leader   *Peer
	LastPing time.Time
	Peers    Map[string, *Peer]
}

func (r *RaftHub) AddPeer(name, address string) {
	r.Peers.Put(name, &Peer{
		Address: address,
		Name:    name,
		State:   FOLLOWER,
	})

	go r.ReplicateLogs()
}

func (r RaftHub) GetPeer(name string) *Peer {
	defer common.HandlePanic("get_peer")

	peer := r.Peers.Get(name)
	if peer == nil {
		panic(common.ErrInvalidPeer)
	}

	return peer
}

func (r *RaftHub) RemovePeer(name string) {
	r.Peers.Delete(name)
}

func (r *RaftHub) InitiateElection() {
	defer common.HandlePanic("initiate_election")

	log.Println("Initiating new election")
	electionRequest := &protocol.InitiateElectionRequest{
		Address: r.Self.Address,
		Name:    r.Self.Name,
		Term:    1,
	}

	r.ChangeLeader(r.Self.Name)
	r.Self.State = CANDIDATE
	ctx := context.Background()
	for _, peer := range r.Peers.entry {
		peer.Send(ctx, electionRequest, ElectionTimeout())
	}

	go r.ReplicateLogs()
}

func (r *RaftHub) CheckLeaderHealth() {
	ticker := time.NewTicker(common.LeaderHealthCheckDelay)
	defer ticker.Stop()

	for range ticker.C {
		if r.Self != r.Leader {
			if !r.Leader.CheckIsAlive() {
				r.Leader.State = DEAD
				r.InitiateElection()

				go r.ReplicateLogs()
			}
		}
	}
}

func (r *RaftHub) ChangeLeader(name string) {
	peer := r.GetPeer(name)
	if peer == nil {
		panic(common.ErrInvalidLeader)
	}

	// test Leader connection before making as Leader
	peer.State = LEADER
	r.Leader = peer

	go r.ReplicateLogs()
}

func (r *RaftHub) CheckFollowersHealth() {
	ticker := time.NewTicker(common.FollowerHealthCheckDelay)
	defer ticker.Stop()

	for range ticker.C {

		for _, peer := range r.Peers.entry {
			if peer.State != LEADER {
				if !peer.CheckIsAlive() {
					peer.State = DEAD

					go r.ReplicateLogs()
				}
			}
		}
	}
}

func (r *RaftHub) SendConnectionRequest(ctx context.Context, request *protocol.AddPeerRequest) error {
	defer common.HandlePanic("request_connection")

	conn, err := grpc.Dial(r.Leader.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		peer := r.GetPeer(request.Peer.Name)
		peer.State = DEAD
		return errors.New(common.ErrPeerUnavailable.Error())
	}

	client := protocol.NewClusterClient(conn)

	resp, _ := client.AddPeer(context.Background(), request)
	for _, peerStruct := range resp.Peers {
		if peerStruct.State != string(LEADER) {
			r.AddPeer(peerStruct.Name, peerStruct.Address)
		}
	}

	return nil
}

func (r *RaftHub) GetLogFilename(ctx context.Context) (*protocol.RPCResponse, error) {
	conn, err := grpc.Dial(r.Leader.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New(common.ErrPeerUnavailable.Error())
	}

	client := protocol.NewClusterClient(conn)
	return client.GetLogFilename(ctx, nil)
}

func (r *RaftHub) Vote(name string) {
	// TODO : Implementations
	// initiate election with term number & request vote from other Peers
	// if term number is lesser that the received number, change to a follower
	// randomize election timeout

	for _, peer := range r.Peers.entry {
		if peer.State == LEADER {
			continue
		}

		if !peer.CheckIsAlive() {
			peer.State = DEAD
		}
	}
}

func (r *RaftHub) ReplicateLogs() {
	peerStates := make([]*protocol.Peer, 0)
	for _, state := range r.Peers.GetEntries() {
		peerStates = append(peerStates, &protocol.Peer{
			Name:    state.Name,
			Address: state.Address,
			State:   string(state.State),
		})
	}

	for _, node := range r.Peers.GetEntries() {
		node.ReplicateLogs(context.Background(), &protocol.PeerList{
			Peers: peerStates,
		})
	}
}

func ElectionTimeout() time.Duration {
	// return time.Duration((time.Duration(rand.Intn(300-500+1)) + 300) * time.Millisecond)
	return 3 * time.Second
}
