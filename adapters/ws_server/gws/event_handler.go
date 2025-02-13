package _gws_lib

import (
	_logger "github.com/dacalin/ws_gateway/logger"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	"github.com/lxzan/gws"
	"time"
)

const (
	PingWait = 10 * time.Second
)

type EventHandler struct {
	hub 		  _iserver.Hub
	pingInterval   time.Duration
	fnOnConnect    _iserver.FnOnConnect
	fnOnDisconnect _iserver.FnOnDisconnect
	fnOnPing       _iserver.FnOnPing
	fnOnMessage    _iserver.FnOnMessage
}

// get connection id from socket session
func getCid(socket *gws.Conn) _connection_id.ConnectionId {
	cid, _ := socket.Session().Load("cid")
	return cid.(_connection_id.ConnectionId)
}

// OnOpen is called when a new connection is opened.
func (e *EventHandler) OnOpen(socket *gws.Conn) {
	_logger.Instance().Printf("OnOpen, cid=%s,", getCid(socket))

	_ = socket.SetDeadline(time.Now().Add(e.pingInterval + PingWait))

	conn := CreateClientConnection(socket)
	e.hub.Set(conn.ConnectionId(), conn)

	if e.fnOnConnect != nil {
		paramsI, _ := socket.Session().Load("params")
		params := paramsI.(map[string]string)
		e.fnOnConnect(conn.ConnectionId(), params)
	}

}

// OnClose is called when a connection is closed.
func (e *EventHandler) OnClose(socket *gws.Conn, err error) {
	_logger.Instance().Printf("onclose, cid=%s, msg=%s\n", getCid(socket), err.Error())

	connId := getCid(socket)

	if e.fnOnDisconnect != nil {
		e.fnOnDisconnect(connId)
	}

	e.hub.Delete(connId)
}

// OnPing is called when a ping message is received.
func (e *EventHandler) OnPing(socket *gws.Conn, payload []byte) {
	// Update the deadline for the connection
	_ = socket.SetDeadline(time.Now().Add(e.pingInterval + PingWait))

	connId := getCid(socket)

	if e.fnOnPing != nil {
		e.fnOnPing(connId)
	}

	e.hub.Send(connId, []byte("pong"))
}

func (e *EventHandler) OnPong(socket *gws.Conn, payload []byte) {

}

// OnMessage is called when a message is received.
func (e *EventHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	// Update the deadline on message
	_ = socket.SetDeadline(time.Now().Add(e.pingInterval + PingWait))

	_logger.Instance().Printf("OnMessage, cid=%s, msg=%s\n", getCid(socket), string(message.Bytes()))

	defer message.Close()

	connId := getCid(socket)

	// chrome websocket
	if b := message.Data.Bytes(); len(b) == 4 && string(b) == "ping" {
		e.OnPing(socket, nil)
		return
	}

	if e.fnOnMessage != nil {
		e.fnOnMessage(connId, message.Bytes())
	}
}
