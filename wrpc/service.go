package wrpc

import (
	"context"
)

type ID uint32

// Service represents a wRPC service.
// All functions are thread-safe.
type Service interface {
	// ID returns the svc_id for a wRPC service.
	ID() ID
	// RPC returns a handler for a rpc_id.
	// Nil handler must be returned if the rpc_id is invalid.
	RPC(ID) RPC
}

// RPC represent a wRPC RPC defined for a Service.
type RPC interface {
	// Handler returns the handler that will be invoked to complete an RPC.
	Handler() Handler
}

// Handler is the handler for a specific wRPC call.
// decoder function is used to unmarshal the request.
// The handler must return the response of the interface and any error.
// Error must not be used to signal application-level failures.
type Handler func(ctx context.Context, peer *Peer, decoder func(interface{}) error) (
	interface{}, error)

// ServiceImpl implements Service
type ServiceImpl struct {
	SvcID ID
	RPCs  map[ID]RPC
}

type RPCImpl struct {
	HandlerFunc Handler
}

func (r RPCImpl) Handler() Handler {
	return r.HandlerFunc
}
