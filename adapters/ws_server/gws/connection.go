package _gws_lib

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iconnection "github.com/dacalin/ws_gateway/ports/connection"
	"github.com/lxzan/gws"
	"sync"
)

var _ _iconnection.Connection = (*ClientConnection)(nil)

type ClientConnection struct {
	_iconnection.Connection
	socket    *gws.Conn
	cid       _connection_id.ConnectionId
	sendMutex sync.Mutex
}

func getCID(socket *gws.Conn) _connection_id.ConnectionId {
	cid, _ := socket.Session().Load("cid")
	return cid.(_connection_id.ConnectionId)
}

func CreateClientConnection(socket *gws.Conn) *ClientConnection {

	return &ClientConnection{
		cid:       getCID(socket),
		socket:    socket,
		sendMutex: sync.Mutex{},
	}
}

func (self *ClientConnection) Send(data []byte) {
	self.sendMutex.Lock()
	defer self.sendMutex.Unlock()
	self.socket.WriteMessage(gws.OpcodeText, data)
}

func (self *ClientConnection) ConnectionId() _connection_id.ConnectionId {
	return self.cid
}
