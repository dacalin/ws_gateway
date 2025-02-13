package _gws_lib

import (
	_logger "github.com/dacalin/ws_gateway/logger"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iconnection "github.com/dacalin/ws_gateway/ports/connection"
	"github.com/lxzan/gws"
	"sync"
)

var _ _iconnection.Connection = (*ClientConnection)(nil)

// ClientConnection represents a connection to a client.
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

// CreateClientConnection creates a new client connection.
func CreateClientConnection(socket *gws.Conn) *ClientConnection {

	return &ClientConnection{
		cid:       getCID(socket),
		socket:    socket,
		sendMutex: sync.Mutex{},
	}
}

// Send sends the given data to the client.
func (conn *ClientConnection) Send(data []byte) {
	conn.sendMutex.Lock()
	defer conn.sendMutex.Unlock()
	_logger.Instance().Printf("Connection Send, data=%v\n", data)

	conn.socket.WriteMessage(gws.OpcodeText, data)
}

// ConnectionId Returns the connection id.
func (conn *ClientConnection) ConnectionId() _connection_id.ConnectionId {
	return conn.cid
}
