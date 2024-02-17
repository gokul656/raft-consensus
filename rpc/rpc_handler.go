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
			resp, _ := s.Register(ctx, request.GetAddPeer())
			conn.Send(RPCCollection_Response(resp))
		case *protocol.RPCRequest_RemovePeer:
			resp, _ := s.DeRegister(ctx, request.GetAddPeer())
			conn.Send(RPCCollection_Response(resp))
		case *protocol.RPCRequest_ElectionRequest:
			resp, _ := s.InitiateElection(ctx, request.GetElectionRequest())
			_ = resp
			// conn.Send(RPCStruct_Response())
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

func (s *GRPCServer) Register(ctx context.Context, req *protocol.AddPeerRequest) (*protocol.PeerList, error) {
	newPeer := req.Peer
	if s.hub.GetPeer(newPeer.Name) != nil {
		log.Println("peer alive:", newPeer.Address)
	} else {
		log.Println("adding peer:", newPeer.Address)
	}

	s.hub.AddPeer(newPeer.Name, newPeer.Address)
	return s.GetPeers(ctx, &emptypb.Empty{})
}

func (s *GRPCServer) DeRegister(ctx context.Context, req *protocol.AddPeerRequest) (*protocol.PeerList, error) {
	s.hub.RemovePeer(req.Peer.Name)
	return s.GetPeers(ctx, &emptypb.Empty{})
}

func (s *GRPCServer) NotifyAll(ctx context.Context, req *protocol.Event) (*protocol.PeerList, error) {
	hub := s.hub
	switch req.GetMessage().(type) {
	case *protocol.Event_PeerAddedEvent:
		hub.AddPeer(req.GetPeerAddedEvent().Name, req.GetPeerAddedEvent().Address)
	case *protocol.Event_PeerRemovedEvent:
		hub.RemovePeer(req.GetPeerRemovedEvent().Name)
	case *protocol.Event_PeerStateChangeEvent:
		event := req.GetPeerStateChangeEvent()
		hub.UpdatePeerStatus(event.Name, *event.PeerState.Enum())
		log.Println("Updating event", event.Name, event.PeerState.Enum())
	default:
		log.Println("Updating event", req)
	}

	return s.GetPeers(ctx, nil)
}

func (s *GRPCServer) InitiateElection(ctx context.Context, req *protocol.InitiateElectionMessage) (*protocol.InitiateElectionMessage, error) {
	log.Println("incoming eleciton request")
	return req, nil
}

func (s *GRPCServer) GetPeers(ctx context.Context, _ *emptypb.Empty) (*protocol.PeerList, error) {
	peers := make([]*protocol.Peer, 0)

	peerMap := s.hub.GetPeerList()
	for _, peer := range peerMap.GetEntries() {
		peers = append(peers, &protocol.Peer{Name: peer.Name, Address: peer.Address, State: peer.State.String()})
	}

	return &protocol.PeerList{Peers: peers}, nil
}

func (s *GRPCServer) GetLogFilename(ctx context.Context, _ *emptypb.Empty) (*protocol.RPCResponse, error) {
	return RPCText_Response(common.GetLogfileName(config.GetEnv().LogDir)), nil
}

func NewgRPCServer() protocol.ClusterServer {
	return &GRPCServer{
		hub: internal.Cluster,
	}
}
