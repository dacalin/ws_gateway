package _gws_lib

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iconnection "github.com/dacalin/ws_gateway/ports/connection"
	"github.com/lxzan/gws"
)

type ClientConnection struct {
	_iconnection.Connection
	socket *gws.Conn
	cid    _connection_id.ConnectionId
}

func getCID(socket *gws.Conn) _connection_id.ConnectionId {
	cid, _ := socket.Session().Load("sid")
	return cid
}

func CreateClientConnection(socket *gws.Conn) *ClientConnection {

	return &ClientConnection{
		cid:    getCID(socket),
		socket: socket,
	}
}

func (self *ClientConnection) Write(data []byte) {
	self.socket.WriteMessage(gws.OpcodeText, data)
}

func (self *ClientConnection) ConnectionId() _connection_id.ConnectionId {
	return self.cid
}
