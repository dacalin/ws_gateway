package gateway

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	_ihub "github.com/dacalin/ws_gateway/ports/hub"
	"sync"
)

type groupName string

var _ _igateway.Gateway = (*Gateway)(nil)

type Gateway struct {
	_igateway.Gateway
	groups map[groupName]*ConnectionMap
	hub    _ihub.Hub
}

var lock = &sync.Mutex{}
var instance *Gateway

func Instance() *Gateway {
	return instance
}

func New(hub _ihub.Hub) *Gateway {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		if instance == nil {
			instance = &Gateway{
				groups: make(map[groupName]*ConnectionMap),
				hub:    hub,
			}
		}
	}
	return instance
}

func (self *Gateway) Send(cid _connection_id.ConnectionId, data []byte) {
	self.hub.Send(cid, data)
}

func (self *Gateway) Broadcast(group string, data []byte) {
	items := self.groups[groupName(group)].Items()

	for conn, _ := range items {
		self.Send(conn, data)
	}
}

func (self *Gateway) SetGroup(cid _connection_id.ConnectionId, group string) {
	if self.groups[groupName(group)] == nil {
		newMap := NewConnectionMap()
		self.groups[groupName(group)] = &newMap
	}

	self.groups[groupName(group)].Set(cid)
}

func (self *Gateway) RemoveGroup(cid _connection_id.ConnectionId, group string) {
	self.groups[groupName(group)].Delete(cid)
}
