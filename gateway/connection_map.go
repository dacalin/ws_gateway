package gateway

import _iconnection "github.com/dacalin/ws_gateway/models/connection_id"

type ConnectionMap struct {
	items map[_iconnection.ConnectionId]string
}

func NewConnectionMap() ConnectionMap {
	return ConnectionMap{
		items: make(map[_iconnection.ConnectionId]string),
	}
}

func (self *ConnectionMap) Set(cid _iconnection.ConnectionId) {
	self.items[cid] = ""
}

func (self *ConnectionMap) Items() map[_iconnection.ConnectionId] string {
	return self.items
}

func (self *ConnectionMap) Delete(cid _iconnection.ConnectionId) {
	delete(self.items, cid)
}
