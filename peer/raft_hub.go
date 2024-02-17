package peer

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/protocol"
)

type RaftHub struct {
	term    int
	Self    *Peer
	Leader  *Peer
	peers   Map[string, *Peer]
	EventCh chan string
}

func (r *RaftHub) AddPeer(name, address string) {
	r.peers.Put(name, &Peer{
		Address: address,
		Name:    name,
		State:   *protocol.PeerState_FOLLOWER.Enum(),
	})

	go r.notifyAll(&protocol.Event{Message: &protocol.Event_PeerAddedEvent{
		PeerAddedEvent: &protocol.PeerAddedEvent{Name: name, Address: address, State: protocol.PeerState_FOLLOWER},
	}})
}

func (r *RaftHub) RemovePeer(name string) {
	r.peers.Delete(name)
	go r.notifyAll(&protocol.Event{Message: &protocol.Event_PeerRemovedEvent{
		PeerRemovedEvent: &protocol.PeerRemovedEvent{Name: name},
	}})
}

func (r *RaftHub) InitiateElection() {
	electionRequest := &protocol.InitiateElectionMessage{
		Name:    r.Self.Name,
		Address: r.Self.Address,
		Term:    uint64(r.term),
	}

	log.Println("[INFO] Initiating election", r.peers)
	for _, peer := range r.peers.GetEntries() {
		if peer != r.Self && peer.State != *protocol.PeerState_DEAD.Enum() {
			response, err := r.InvokePeerElection(context.Background(), peer.Address, electionRequest, ElectionTimeout())
			if err == nil {
				log.Println(response.String())
			}
		}
	}
}

func (r *RaftHub) notifyAll(event *protocol.Event) {
	if r.Self == r.Leader {
		for _, peer := range r.peers.GetEntries() {
			if peer != r.Self {
				_ = peer.SendMesagge(event)
				// TODO : Handle error
			}
		}
	}
}

func (r RaftHub) GetPeer(name string) *Peer {
	return r.peers.Get(name)
}

func (r RaftHub) GetPeerList() Map[string, *Peer] {
	return r.peers
}

func (r *RaftHub) UpdatePeerStatus(name string, state protocol.PeerState) {
	peer := r.GetPeer(name)
	if peer == nil {
		log.Println("[INFO] Peer not registered", name)
		return
	}

	if peer.State.Enum() != state.Enum() {
		peer.State = state
		go r.notifyAll(&protocol.Event{
			Message: &protocol.Event_PeerStateChangeEvent{
				PeerStateChangeEvent: &protocol.PeerStateChangeEvent{
					Name:      name,
					PeerState: state,
				},
			},
		})
	}
}

func (r *RaftHub) ChangeLeader(name string) error {
	peer := r.GetPeer(name)
	if peer == nil {
		return common.ErrInvalidLeader
	}

	// test Leader connection before making as Leader
	r.Leader = peer
	r.UpdatePeerStatus(name, *protocol.PeerState_LEADER.Enum())
	return nil
}

func (r *RaftHub) CheckLeaderHealth() {
	ticker := time.NewTicker(common.LeaderHealthCheckDelay)
	defer ticker.Stop()

	for range ticker.C {
		if r.Self != r.Leader {
			if !r.Leader.CheckIsAlive() {
				r.UpdatePeerStatus(r.Leader.Name, *protocol.PeerState_DEAD.Enum())
				r.InitiateElection()
			}
		}
	}
}

func (r *RaftHub) CheckFollowersHealth(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for range ticker.C {
		for _, peer := range r.peers.entry {
			if peer.State != *protocol.PeerState_LEADER.Enum() {
				existingPeer := r.GetPeer(peer.Name)
				if !peer.CheckIsAlive() && existingPeer.State != *protocol.PeerState_DEAD.Enum() {
					log.Printf("peer dead: %s", peer.Name)
					r.UpdatePeerStatus(peer.Name, *protocol.PeerState_DEAD.Enum())
				}
			}
		}
	}
}

func (r *RaftHub) Synchronize() {
	log.Println("[INFO] Synchronizing with Leader node...")
	peerList, err := r.GetPeers(context.Background())
	if err != nil {
		log.Fatalln("[CRITICAL] Unable to Sync", err)
	}

	for _, peer := range peerList.GetPeers() {
		peerData := r.GetPeer(peer.Name)
		if peerData == nil {
			r.AddPeer(peer.Name, peer.Address)
		}

		r.UpdatePeerStatus(peer.Name, *protocol.PeerState(protocol.PeerState_value[peer.State]).Enum())
	}

	log.Println("[INFO] Synchronizing success", r.peers)
}

func ElectionTimeout() time.Duration {
	return time.Duration(rand.Intn(301)+200) * time.Millisecond
}

func NewRaft(self *Peer) *RaftHub {
	return &RaftHub{
		Self:    self,
		peers:   NewMap[string, *Peer](),
		EventCh: make(chan string),
	}
}
