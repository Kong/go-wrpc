package wrpc

import (
	"context"
	"testing"
)

func TestRPCImpl_Handler(t *testing.T) {
	t.Run("empty handler returns nil", func(t *testing.T) {
		rpc := RPCImpl{}
		if rpc.Handler() != nil {
			t.Errorf("expected nil handler")
		}
	})
	t.Run("initialized handler returns a handler", func(t *testing.T) {
		rpc := RPCImpl{
			HandlerFunc: func(ctx context.Context, decoder func(interface{}) error) (interface{}, error) {
				return nil, nil
			},
		}
		if rpc.Handler() == nil {
			t.Errorf("expected a handler")
		}
	})
}
