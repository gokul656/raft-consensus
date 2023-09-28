package peer

import (
	"context"
	"time"
)

type PeerState string

const (
	LEADER    PeerState = "LEADER"
	CANDIDATE PeerState = "CANDIDATE"
	FOLLOWER  PeerState = "FOLLOWER"
	DEAD      PeerState = "DEAD"
)

type Peer struct {
	Address string
	Name    string
	State   PeerState
}

func (p *Peer) UpdatePeerState(newState PeerState) {
	p.State = newState
}

func (p *Peer) Send(ctx context.Context, message []byte, timeout time.Duration) ([]byte, error) {
	// TODO : Implementations
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return nil, nil
}

func (p *Peer) CheckIsAlive() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	return p.ping(ctx, p.Address)
}

func (p *Peer) ping(ctx context.Context, address string) bool {
	ticker := time.NewTicker(time.Duration(1) * time.Second)
	select {
	case <-ticker.C:
		// TODO : if status is up return true
		p.Send(ctx, []byte("ping"), ElectionTimeout())
		return true
	case <-ctx.Done():
		ticker.Stop()
		return false
	}
}
