package rpc

import (
	"github.com/gokul656/raft-consensus/protocol"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func RPCText_Response(text string) *protocol.RPCResponse {
	return &protocol.RPCResponse{
		Message: &protocol.RPCResponse_Text{
			Text: text,
		},
	}
}

func RPCStruct_Response(data *structpb.Struct) *protocol.RPCResponse {
	return &protocol.RPCResponse{
		Message: &protocol.RPCResponse_Struct{
			Struct: data,
		},
	}
}

func RPCCollection_Response(dataset []*anypb.Any) *protocol.RPCResponse {
	return &protocol.RPCResponse{
		Message: &protocol.RPCResponse_Collection{
			Collection: &protocol.CollectionMessage{
				Data: dataset,
			},
		},
	}
}
