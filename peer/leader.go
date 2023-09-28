package peer

import (
	"context"
	"math/rand"
	"time"

	"github.com/gokul656/raft-consensus/common"
)

type Vote struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Term    uint   `json:"term"`
}

type RaftHub struct {
	Ctx    context.Context
	Self   *Peer
	Leader *Peer
	Peers  Map[string, *Peer]
}

func (r *RaftHub) BootstrapPeer() {
	if r.Leader == nil {
		r.Vote(r.Self.Name)
		r.Leader = r.Self
	}
}

func (r *RaftHub) AddPeer(name, address string) {
	r.Peers.Put(name, &Peer{
		Address: address,
		Name:    name,
		State:   CANDIDATE,
	})
}

func (r *RaftHub) RemovePeer(name string) {
	r.Peers.Delete(name)
}

func (r *RaftHub) InitiateElection() {
	vote := &Vote{
		Address: r.Self.Address,
		Name:    r.Self.Name,
		Term:    1,
	}

	for _, peer := range r.Peers.entry {
		marshalled, _ := common.ToByte(vote)
		response, _ := peer.Send(r.Ctx, marshalled, ElectionTimeout())
		unmarshalled, _ := common.FromByte[Vote](response, Vote{})

		if unmarshalled.Term > vote.Term {
			r.ChangeLeader(unmarshalled.Name)
			break
		}

		vote.Term++
	}
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

func (r *RaftHub) ChangeLeader(name string) {
	if r.GetPeer(name) == nil {
		panic(common.InvalidLeader)
	}

	// test Leader connection before making as Leader
	r.Leader = r.GetPeer(name)
}

func (r RaftHub) GetPeer(name string) *Peer {
	return r.Peers.Get(&name)
}

func (r *RaftHub) ReplicateLogs() {
	panic("not yet implemented")
}

func (r RaftHub) CheckIsAlive(name string) bool {
	// TODO : Implement ping to peer
	return false
}

func (r *RaftHub) StartHeartBeats() {
	ctx := context.Background()
	ticker := time.NewTicker(time.Duration(3) * time.Second)
	select {
	case <-ticker.C:
		for _, peer := range r.Peers.entry {
			if peer.State == LEADER {
				continue
			}

			peer.Send(ctx, []byte("ping"), ElectionTimeout())
		}
	}
}

func ElectionTimeout() time.Duration {
	return time.Duration((time.Duration(rand.Intn(300-500+1)) + 300) * time.Millisecond)
}

type Map[K any, V any] struct {
	entry map[*K]V
}

func (m *Map[K, V]) Get(key *K) V {
	return m.entry[key]
}

func (m *Map[K, V]) GetEntries() map[*K]V {
	return m.entry
}

func (m *Map[K, V]) Put(key K, value V) {
	m.entry[&key] = value
}

func (m *Map[K, V]) Delete(key K) {
	delete(m.entry, &key)
}

func NewMap[K any, V any]() Map[K, V] {
	return Map[K, V]{entry: map[*K]V{}}
}
