// Code generated by protoc-gen-go-wrpc. DO NOT EDIT
// protoc-gen-go-wrpc version: (devel)

package proto

import (
	context "context"
	wrpc "github.com/kong/go-wrpc/wrpc"
)

type Echo interface {
	Echo(context.Context, *wrpc.Peer, *EchoRPCRequest) (*EchoRPCResponse, error)
}

func PrepareEchoEchoRequest(in *EchoRPCRequest) (wrpc.Request, error) {
	return wrpc.CreateRequest(1, 1, in)
}

type EchoClient struct {
	Peer *wrpc.Peer
}

func (c *EchoClient) Echo(ctx context.Context, in *EchoRPCRequest) (*EchoRPCResponse, error) {
	err := c.Peer.VerifyRPC(1, 1)
	if err != nil {
		return nil, err
	}

	req, err := PrepareEchoEchoRequest(in)
	if err != nil {
		return nil, err
	}

	var out EchoRPCResponse
	err = c.Peer.DoRequest(ctx, req, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

type EchoServer struct {
	Echo Echo
}

func (s *EchoServer) ID() wrpc.ID {
	return 1
}

func (s *EchoServer) RPC(rpc wrpc.ID) wrpc.RPC {
	switch rpc {
	case 1:
		return wrpc.RPCImpl{
			HandlerFunc: func(ctx context.Context, peer *wrpc.Peer, decode func(interface{}) error) (interface{}, error) {
				var in EchoRPCRequest
				err := decode(&in)
				if err != nil {
					return nil, err
				}
				return s.Echo.Echo(ctx, peer, &in)
			},
		}
	default:
		return nil
	}
}
