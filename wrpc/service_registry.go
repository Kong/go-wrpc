package wrpc

import (
	"fmt"
	"sync"
)

type serviceRegistry struct {
	once     sync.Once
	services map[ID]Service
	lock     sync.RWMutex
}

func (r *serviceRegistry) init() {
	r.services = map[ID]Service{}
}

func (r *serviceRegistry) Add(s Service) error {
	r.once.Do(r.init)
	if s == nil {
		return fmt.Errorf("nil service")
	}
	if s.ID() == 0 {
		return fmt.Errorf("invalid service id(0)")
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.services[s.ID()]; ok {
		return fmt.Errorf("service already registered")
	}
	r.services[s.ID()] = s
	return nil
}

func (r *serviceRegistry) Get(svcID, rpcID ID) (RPC, error) {
	r.once.Do(r.init)
	if svcID == 0 {
		return nil, fmt.Errorf("invalid service id(0)")
	}
	if rpcID == 0 {
		return nil, fmt.Errorf("invalid rpc id(0)")
	}

	r.lock.RLock()
	defer r.lock.RUnlock()
	s, ok := r.services[svcID]
	if !ok {
		return nil, fmt.Errorf("invalid service(%d)", svcID)
	}
	handler := s.RPC(rpcID)
	if handler == nil {
		return nil, fmt.Errorf("invalid rpc(%d) for service(%d)", rpcID, svcID)
	}
	return handler, nil
}
