package _gws_lib

import (
	"fmt"
	_hub "github.com/dacalin/ws_gateway/hub"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_server "github.com/dacalin/ws_gateway/ports/server"
	"github.com/lxzan/gws"
	"log"
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
	return cid
}

func (self *EventHandler) OnOpen(socket *gws.Conn) {
	fmt.Println("OnOpen")
	fmt.Println(getCid(socket))

	_ = socket.SetDeadline(time.Now().Add(self.pingInterval + PingWait))
	conn := CreateClientConnection(socket)
	_hub.Instance().Set(conn.ConnectionId(), conn)

	if self.fnOnConnect != nil {
		self.fnOnConnect(conn.ConnectionId(), map[string]string{})
	}

}

func (self *EventHandler) OnClose(socket *gws.Conn, err error) {
	log.Printf("onclose, cid=%s, msg=%s\n", getCid(socket), err.Error())

	connId := getCid(socket)
	_hub.Instance().Delete(connId)

	if self.fnOnDisconnect != nil {
		self.fnOnDisconnect(connId)
	}
}

func (self *EventHandler) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(self.pingInterval + PingWait))

	//_ws_api.OnPing(getCid(socket))
	connId := getCid(socket)

	if self.fnOnPing != nil {
		self.fnOnPing(connId)
	}
}

func (self *EventHandler) OnPong(socket *gws.Conn, payload []byte) {}

func (self *EventHandler) OnMessage(socket *gws.Conn, message *gws.Message) {
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
