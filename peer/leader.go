package peer

import "github.com/gokul656/raft-consensus/common"

type RaftHub struct {
	leader string
	peers  map[string]string
}

func (r *RaftHub) AddPeer(name, address string) {
	r.peers[name] = address
}

func (r *RaftHub) RemovePeer(name string) {
	delete(r.peers, name)
}

func (r *RaftHub) Vote(name string) {
	// TODO : Implementations
	// initiate election with term number & request vote from other peers
	// if term number is lesser that the received number, change to a follower
	// randomize election timeout
}

func (r *RaftHub) ChangeLeader(name string) {
	if r.GetPeer(name) == "" {
		panic(common.INVALID_LEADER)
	}

	// test leader connection before make as leader
	r.leader = name
}

func (r RaftHub) GetPeer(name string) string {
	return r.peers[name]
}

func (r RaftHub) CheckIsAlive(name string) bool {
	// TODO : Implement ping to peer
	return false
}
