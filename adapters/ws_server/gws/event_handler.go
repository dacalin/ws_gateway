package _gws_lib

import (
	_gws_hub "github.com/dacalin/ws_gateway/adapters/ws_server/gws/hub"
	_logger "github.com/dacalin/ws_gateway/logger"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_server "github.com/dacalin/ws_gateway/ports/server"
	"github.com/lxzan/gws"
	"time"
)

const (
	PingWait = 10 * time.Second
)

type EventHandler struct {
	pingInterval   time.Duration
	fnOnConnect    _server.FnOnConnect
	fnOnDisconnect _server.FnOnDisconnect
	fnOnPing       _server.FnOnPing
	fnOnMessage    _server.FnOnMessage
}

func getCid(socket *gws.Conn) _connection_id.ConnectionId {
	cid, _ := socket.Session().Load("cid")
	return cid.(_connection_id.ConnectionId)
}

func (self *EventHandler) OnOpen(socket *gws.Conn) {
	_logger.Instance().Printf("OnOpen, cid=%s,", getCid(socket))

	_ = socket.SetDeadline(time.Now().Add(self.pingInterval + PingWait))

	conn := CreateClientConnection(socket)
	_gws_hub.Instance().Set(conn.ConnectionId(), conn)

	if self.fnOnConnect != nil {
		paramsI, _ := socket.Session().Load("params")
		params := paramsI.(map[string]string)
		self.fnOnConnect(conn.ConnectionId(), params)
	}

}

func (self *EventHandler) OnClose(socket *gws.Conn, err error) {
	_logger.Instance().Printf("onclose, cid=%s, msg=%s\n", getCid(socket), err.Error())

	connId := getCid(socket)

	if self.fnOnDisconnect != nil {
		self.fnOnDisconnect(connId)
	}

	_gws_hub.Instance().Delete(connId)
}

func (self *EventHandler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(self.pingInterval + PingWait))
	connId := getCid(socket)

	if self.fnOnPing != nil {
		self.fnOnPing(connId)
	}

	_gws_hub.Instance().Send(connId, []byte("pong"))
}

func (self *EventHandler) OnPong(socket *gws.Conn, payload []byte) {

}

func (self *EventHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
	_ = socket.SetDeadline(time.Now().Add(self.pingInterval + PingWait))

	_logger.Instance().Printf("OnMessage, cid=%s, msg=%s\n", getCid(socket), string(message.Bytes()))

	defer message.Close()

	connId := getCid(socket)

	// chrome websocket
	if b := message.Data.Bytes(); len(b) == 4 && string(b) == "ping" {
		self.OnPing(socket, nil)
		return
	}

	if self.fnOnMessage != nil {
		self.fnOnMessage(connId, message.Bytes())
	}
}
