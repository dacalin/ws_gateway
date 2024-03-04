package gateway

import (
	_logger "github.com/dacalin/ws_gateway/logger"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	_ihub "github.com/dacalin/ws_gateway/ports/hub"
	"sync"
)

type groupName string

var _ _igateway.Gateway = (*Gateway)(nil)

type Gateway struct {
	_igateway.Gateway
	groups sync.Map //groups map[groupName]ConnectionMap
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
		instance = &Gateway{
			groups: sync.Map{},
			hub:    hub,
		}
	}
	return instance
}

func (self *Gateway) Send(cid _connection_id.ConnectionId, data []byte) {
	self.hub.Send(cid, data)
}

func (self *Gateway) Broadcast(group string, data []byte) {
	_logger.Instance().Printf("Broadcast, group=%s  data=%s", group, data)

	connMapI, _ := self.groups.Load(groupName(group))
	connMap := connMapI.(ConnectionMap)

	for cid, _ := range connMap.Items() {
		self.Send(cid, data)
	}
}

func (self *Gateway) SetGroup(cid _connection_id.ConnectionId, group string) {
	connMapI, found := self.groups.Load(groupName(group))
	var connMap ConnectionMap

	if found == false {
		newMap := NewConnectionMap()
		self.groups.Store(groupName(group), newMap)
		connMap = newMap

	} else {
		connMap = connMapI.(ConnectionMap)
	}

	connMap.Set(cid)
}

func (self *Gateway) RemoveGroup(cid _connection_id.ConnectionId, group string) {
	connMapI, found := self.groups.Load(groupName(group))
	var connMap ConnectionMap

	if found == true {
		connMap = connMapI.(ConnectionMap)
		connMap.Delete(cid)

		if len(connMap.items) == 0 {
			self.groups.Delete(groupName(group))
		}
	}

}
