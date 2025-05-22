package gateway

import (
	"sync"

	_iconnection "github.com/dacalin/ws_gateway/models/connection_id"
)

type ConnectionMap struct {
	items map[_iconnection.ConnectionId]string
	mu    sync.RWMutex
}

func NewConnectionMap() ConnectionMap {
	return ConnectionMap{
		items: make(map[_iconnection.ConnectionId]string),
		mu:    sync.RWMutex{},
	}
}

func (self *ConnectionMap) Set(cid _iconnection.ConnectionId) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.items[cid] = ""
}

func (self *ConnectionMap) Items() map[_iconnection.ConnectionId]string {
	self.mu.RLock()
	defer self.mu.RUnlock()
	copiedItems := make(map[_iconnection.ConnectionId]string)
	for key, value := range self.items {
		copiedItems[key] = value
	}
	return copiedItems
}

func (self *ConnectionMap) Delete(cid _iconnection.ConnectionId) {
	self.mu.Lock()
	defer self.mu.Unlock()
	delete(self.items, cid)
}

func (self *ConnectionMap) Len() int {
	self.mu.RLock()
	defer self.mu.RUnlock()
	return len(self.items)
}
