package peer

import (
	"context"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Peer struct {
	Address string
	Name    string
	State   protocol.PeerState `protobuf:"enum=State" json:"State"`
}

func (p *Peer) CheckIsAlive() bool {
	client, err := GetClientConnection(p.Address)
	if err != nil {
		return false
	}
	_, err = client.Ping(context.Background(), &emptypb.Empty{})
	return err == nil
}

func (p *Peer) SendMesagge(event *protocol.Event) error {
	client, err := GetClientConnection(p.Address)
	if err != nil {
		return common.ErrPeerUnavailable
	}

	_, err = client.NotifyAll(context.Background(), event)
	if err != nil {
		return common.ErrPeerUnavailable
	}

	return nil
}
