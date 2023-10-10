package peer

import (
	"context"
	"log"
	"time"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func GetClientConnection(address string) (protocol.ClusterClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, common.ErrPeerUnavailable
	}
	client := protocol.NewClusterClient(conn)

	return client, nil
}

func (r *RaftHub) Register(ctx context.Context, request *protocol.AddPeerRequest) error {
	defer common.HandlePanic("request_connection")

	log.Println("[INFO] Registering as FOLLOWER", r.Self.Name)
	client, err := GetClientConnection(r.Leader.Address)
	if err != nil {
		return common.ErrPeerUnavailable
	}

	_, err = client.Register(ctx, request)
	return err
}

func (r *RaftHub) DeRegister(ctx context.Context, request *protocol.AddPeerRequest) error {
	client, err := GetClientConnection(r.Leader.Address)
	if err != nil {
		return common.ErrPeerUnavailable
	}

	_, err = client.DeRegister(ctx, request)
	return err
}

func (r *RaftHub) InvokePeerElection(ctx context.Context, address string, message *protocol.InitiateElectionMessage, timeout time.Duration) (*protocol.InitiateElectionMessage, error) {
	client, err := GetClientConnection(address)
	if err != nil {
		return nil, common.ErrPeerUnavailable
	}

	response, err := client.InitiateElection(ctx, message)
	if err != nil {
		return nil, common.ErrPeerUnavailable
	}

	return response, nil
}

func (r *RaftHub) GetLogFilename(ctx context.Context) (*protocol.RPCResponse, error) {
	client, err := GetClientConnection(r.Leader.Address)
	if err != nil {
		return nil, common.ErrPeerUnavailable
	}

	return client.GetLogFilename(ctx, nil)
}

func (r *RaftHub) GetPeers(ctx context.Context) (*protocol.PeerList, error) {
	client, err := GetClientConnection(r.Leader.Address)
	if err != nil {
		return nil, common.ErrPeerUnavailable
	}

	peerList, err := client.GetPeers(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	return peerList, nil
}
