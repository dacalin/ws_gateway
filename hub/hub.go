package _hub

import (
	_ihub "github.com/dacalin/ws_gateway/ports/hub"
	"sync"
)

var _ _ihub.Hub = (*Hub)(nil)

type Hub struct {
	_ihub.Hub
}

var lock = &sync.Mutex{}
var instance _ihub.Hub

func Instance() _ihub.Hub {
	return instance
}

func New(hub _ihub.Hub) _ihub.Hub {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = hub
	}

	return instance
}
