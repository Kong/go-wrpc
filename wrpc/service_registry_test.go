package wrpc

import (
	"testing"
)

type serviceImpl struct {
	id   ID
	rpcs map[ID]RPC
}

func (s serviceImpl) ID() ID {
	return s.id
}

func (s serviceImpl) RPC(id ID) RPC {
	return s.rpcs[id]
}

func Test_serviceRegistry_Add(t *testing.T) {
	t.Run("adding service with id 0 errors out", func(t *testing.T) {
		reg := &serviceRegistry{}
		err := reg.Add(serviceImpl{
			id: 0,
		})
		if err == nil {
			t.Errorf("expected an error")
		}
	})
	t.Run("adding service with id > 0 works", func(t *testing.T) {
		reg := &serviceRegistry{}
		err := reg.Add(serviceImpl{
			id: 1,
		})
		if err != nil {
			t.Errorf("didn't expect an error")
		}
	})
	t.Run("re-adding a service errors out", func(t *testing.T) {
		reg := &serviceRegistry{}
		err := reg.Add(serviceImpl{
			id: 1,
		})
		if err != nil {
			t.Errorf("didn't expect an error")
		}
		err = reg.Add(serviceImpl{
			id: 1,
		})
		if err == nil {
			t.Errorf("expected an error on re-adding")
		}
	})
	t.Run("adding nil service errors out", func(t *testing.T) {
		reg := &serviceRegistry{}
		err := reg.Add(nil)
		if err == nil {
			t.Errorf("expected an error")
		}
	})
}

func Test_serviceRegistry_Get(t *testing.T) {
	reg := &serviceRegistry{}
	reg.Add(serviceImpl{
		id: 1,
		rpcs: map[ID]RPC{
			42: RPCImpl{},
		},
	})
	t.Run("errors out with svc_id 0 ", func(t *testing.T) {
		_, err := reg.Get(0, 42)
		if err == nil {
			t.Errorf("expected an error")
		}
	})
	t.Run("errors out with rpc_id 0 ", func(t *testing.T) {
		_, err := reg.Get(1, 0)
		if err == nil {
			t.Errorf("expected an error")
		}
	})
	t.Run("errors out with non-existent rpc_id", func(t *testing.T) {
		_, err := reg.Get(1, 43)
		if err == nil {
			t.Errorf("expected an error")
		}
	})
	t.Run("errors out with non-existent svc_id", func(t *testing.T) {
		_, err := reg.Get(2, 42)
		if err == nil {
			t.Errorf("expected an error")
		}
	})
	t.Run("fetches a handler with correct ids", func(t *testing.T) {
		_, err := reg.Get(1, 42)
		if err != nil {
			t.Errorf("didn't expect an error")
		}
	})
}
