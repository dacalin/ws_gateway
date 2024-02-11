package gateway

import (
	_hub "github.com/dacalin/ws_gateway/hub"
	_iconnection "github.com/dacalin/ws_gateway/models/connection_id"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	"sync"
)

type groupName string

type Gateway struct {
	_igateway.Gateway
	groups map[groupName]*ConnectionMap
}

var lock = &sync.Mutex{}
var instance *Gateway

func Instance() *Gateway {
	return instance
}

func New() *Gateway {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		if instance == nil {
			instance = &Gateway{
				groups: make(map[groupName]*ConnectionMap),
			}
		}
	}
	return instance
}
func (self *Gateway) Send(cid _iconnection.ConnectionId, data []byte) {
	conn, found := _hub.Instance().Get(cid)

	if found == false {
		_hub.Instance().PubSub().Publish(cid.Value(), data)
	} else {
		conn.Send(data)
	}

}

func (self *Gateway) Broadcast(group string, data []byte) {
	items := self.groups[groupName(group)].Items()

	for conn, _ := range items {
		self.Send(conn, data)
	}
}

func (self *Gateway) SetGroup(cid _iconnection.ConnectionId, group string) {
	if self.groups[groupName(group)] == nil {
		newMap := NewConnectionMap()
		self.groups[groupName(group)] = &newMap
	}

	self.groups[groupName(group)].Set(cid)
}

func (self *Gateway) RemoveGroup(cid _iconnection.ConnectionId, group string) {
	self.groups[groupName(group)].Delete(cid)
}
