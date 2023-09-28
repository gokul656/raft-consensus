package rpc

import (
	"context"
	"fmt"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
	"github.com/gokul656/raft-consensus/peer"
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type gRPCServer struct {
	hub *peer.RaftHub
	protocol.UnimplementedClusterServer
}

func (s *gRPCServer) ConnectionStream(conn protocol.Cluster_ConnectionStreamServer) error {
	defer common.HandlePanic()

	ctx := context.Background()
	for {
		request, err := conn.Recv()
		if err != nil {
			panic(err)
		}

		switch request.Message.(type) {
		case *protocol.RPCRequest_AddPeer:
			resp, _ := s.AddPeer(ctx, request.GetAddPeer())
			conn.Send(resp)
		case *protocol.RPCRequest_ElectionRequest:
			resp, _ := s.InitiateElection(ctx, request.GetElectionRequest())
			conn.Send(resp)
		case *protocol.RPCRequest_GetPeerList:
			resp, _ := s.GetPeers(ctx, nil)
			conn.Send(resp)
		default:
			defaultMessage := &protocol.RPCResponse_Text{
				Text: "unknown method",
			}
			conn.Send(&protocol.RPCResponse{
				Message: defaultMessage,
			})
		}
	}
}

func (s *gRPCServer) AddPeer(ctx context.Context, req *protocol.AddPeerRequest) (*protocol.RPCResponse, error) {
	s.hub.Peers.Put(req.Name, &peer.Peer{
		Address: req.Address,
		Name:    req.Name,
		State:   peer.FOLLOWER,
	})

	return RPCText_Response("success"), nil
}

func (s *gRPCServer) InitiateElection(ctx context.Context, req *protocol.InitiateElectionRequest) (*protocol.RPCResponse, error) {
	return RPCStruct_Response(nil), nil
}

func (s *gRPCServer) GetPeers(ctx context.Context, _ *emptypb.Empty) (*protocol.RPCResponse, error) {
	peers := make([]*anypb.Any, 0)
	for _, peer := range s.hub.Peers.GetEntries() {
		anyData, _ := anypb.New(&protocol.AddPeerRequest{
			Name:    peer.Name,
			Address: peer.Address,
		})
		peers = append(peers, anyData)
	}

	return RPCCollection_Response(peers), nil
}

func NewgRPCServer() protocol.ClusterServer {
	raftHub := &peer.RaftHub{
		Ctx:   context.Background(),
		Peers: peer.NewMap[string, *peer.Peer](),
	}

	currentInstance := &peer.Peer{
		Address: fmt.Sprintf("localhost:%s", config.GetEnv().RPCPort),
		Name:    config.GetEnv().InstanceID,
		State:   peer.LEADER,
	}
	raftHub.Leader = currentInstance
	raftHub.Peers.Put(currentInstance.Name, currentInstance)

	return &gRPCServer{
		hub: raftHub,
	}
}
