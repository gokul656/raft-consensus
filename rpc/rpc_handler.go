package rpc

import (
	"context"
	"log"

	"github.com/gokul656/raft-consensus/common"
	"github.com/gokul656/raft-consensus/config"
	"github.com/gokul656/raft-consensus/internal"
	"github.com/gokul656/raft-consensus/peer"
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	hub *peer.RaftHub
	protocol.UnimplementedClusterServer
}

func (s *GRPCServer) ConnectionStream(conn protocol.Cluster_ConnectionStreamServer) error {
	defer common.HandlePanic("connection_stream")

	ctx := context.Background()
	for {
		request, err := conn.Recv()
		if err != nil {
			panic(err)
		}

		switch request.Message.(type) {
		case *protocol.RPCRequest_AddPeer:
			resp, _ := s.AddPeer(ctx, request.GetAddPeer())
			conn.Send(RPCCollection_Response(resp))
		case *protocol.RPCRequest_ElectionRequest:
			resp, _ := s.InitiateElection(ctx, request.GetElectionRequest())
			conn.Send(resp)
		case *protocol.RPCRequest_GetPeerList:
			resp, _ := s.GetPeers(ctx, nil)
			conn.Send(RPCCollection_Response(resp))
		default:
			defaultMessage := RPCText_Response("unknown method or method not implemented")
			conn.Send(defaultMessage)
		}
	}
}

func (s *GRPCServer) Ping(ctx context.Context, _ *emptypb.Empty) (*protocol.RPCResponse, error) {
	return RPCText_Response("pong"), nil
}

func (s *GRPCServer) AddPeer(ctx context.Context, req *protocol.AddPeerRequest) (*protocol.PeerList, error) {
	newPeer := req.Peer
	if s.hub.GetPeer(newPeer.Name) != nil {
		log.Println("peer alive:", newPeer.Address)
	} else {
		log.Println("adding peer:", newPeer.Address)
	}

	s.hub.Peers.Put(newPeer.Name, &peer.Peer{
		Address: newPeer.Address,
		Name:    newPeer.Name,
		State:   peer.FOLLOWER,
	})
	return s.GetPeers(ctx, &emptypb.Empty{})
}

func (s *GRPCServer) ReplicateLogs(ctx context.Context, req *protocol.PeerList) (*protocol.PeerList, error) {
	for _, newPeer := range req.Peers {
		existingPeer := s.hub.GetPeer(newPeer.Name)
		if existingPeer == nil {
			s.hub.AddPeer(newPeer.Name, newPeer.Address)
		} else {
			existingPeer.State = peer.States[newPeer.State]
		}

		log.Println("updating peer state:", newPeer.Name, newPeer.State, newPeer.Address)
	}

	log.Println("------------------------------------")
	return s.GetPeers(ctx, nil)
}

func (s *GRPCServer) InitiateElection(ctx context.Context, req *protocol.InitiateElectionRequest) (*protocol.RPCResponse, error) {
	return RPCStruct_Response(nil), nil
}

func (s *GRPCServer) GetPeers(ctx context.Context, _ *emptypb.Empty) (*protocol.PeerList, error) {
	peers := &protocol.PeerList{
		Peers: make([]*protocol.Peer, 0),
	}
	for _, peer := range s.hub.Peers.GetEntries() {
		peerData := &protocol.Peer{
			Name:    peer.Name,
			Address: peer.Address,
			State:   string(peer.State),
		}
		peers.Peers = append(peers.Peers, peerData)
	}

	return peers, nil
}

func (s *GRPCServer) GetLogFilename(ctx context.Context, _ *emptypb.Empty) (*protocol.RPCResponse, error) {
	return RPCText_Response(common.GetLogfileName(config.GetEnv().LogDir)), nil
}

func NewgRPCServer() protocol.ClusterServer {
	return &GRPCServer{
		hub: internal.Cluster,
	}
}
