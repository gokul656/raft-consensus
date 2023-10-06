package peer

import (
	"context"
	"time"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type PeerState string

const (
	LEADER    PeerState = "LEADER"
	CANDIDATE PeerState = "CANDIDATE"
	FOLLOWER  PeerState = "FOLLOWER"
	DEAD      PeerState = "DEAD"
)

var States = map[string]PeerState{
	"LEADER":    LEADER,
	"CANDIDATE": CANDIDATE,
	"FOLLOWER":  FOLLOWER,
	"DEAD":      DEAD,
}

type Peer struct {
	Address string
	Name    string
	State   PeerState
}

func (p *Peer) UpdatePeerState(newState PeerState) {
	p.State = newState
}

func (p *Peer) Send(ctx context.Context, message *protocol.InitiateElectionRequest, timeout time.Duration) ([]byte, error) {
	// TODO : Implementations
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := grpc.Dial(p.Address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := protocol.NewClusterClient(conn)
	_, err = client.InitiateElection(ctx, message)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *Peer) ReplicateLogs(ctx context.Context, peerStates *protocol.PeerList) {
	conn, err := grpc.Dial(p.Address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := protocol.NewClusterClient(conn)
	client.ReplicateLogs(ctx, peerStates)
}

func (p *Peer) CheckIsAlive() bool {
	defer common.HandlePanic("check_is_alive")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	return p.ping(ctx, p.Address)
}

func (p *Peer) ping(ctx context.Context, address string) bool {
	defer common.HandlePanic("ping")

	conn, err := grpc.Dial(p.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(common.CustomError{Err: common.ErrPeerUnavailable, Msg: p.Address})
	}

	client := protocol.NewClusterClient(conn)
	_, err = client.Ping(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(common.CustomError{Err: common.ErrPeerUnavailable, Msg: p.Address})
	}

	return true
}
