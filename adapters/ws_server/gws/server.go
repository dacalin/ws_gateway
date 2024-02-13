package _gws_lib

import (
	_hub "github.com/dacalin/ws_gateway/hub"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	"github.com/lxzan/gws"
	"net/http"
	"strconv"
	"time"
)

type WSServer struct {
	_iserver.Server
	connectionRoute string
	eventHandler    EventHandler
	pubsub          _ipubsub.Client
}

func Create(connectionRoute string, pingInterval int, pubsub _ipubsub.Client) WSServer {
	duration := time.Duration(pingInterval) * time.Second

	eventHandler := EventHandler{
		pingInterval:   duration,
		fnOnConnect:    nil,
		fnOnDisconnect: nil,
		fnOnPing:       nil,
		fnOnMessage:    nil,
	}

	return WSServer{
		eventHandler:    eventHandler,
		connectionRoute: connectionRoute,
		pubsub:          pubsub,
	}
}

func (self *WSServer) Run(port int) {

	_hub.New(self.pubsub)

	upgrader := gws.NewUpgrader(&self.eventHandler, &gws.ServerOption{
		ReadAsyncEnabled: true,         // Parallel messages processing
		CompressEnabled:  true,         // Enable compression
		Recovery:         gws.Recovery, // Exception recovery
	})

	http.HandleFunc("/"+self.connectionRoute, func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			return
		}

		cidParam := request.URL.Query().Get("cid")
		if cidParam == "" {
			return
		}

		cid := _connection_id.New(cidParam)
		socket.Session().Store("cid", cid)

		go func() {
			socket.ReadLoop() // Blocking prevents the context from being GC.
		}()
	})

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func (self *WSServer) OnConnect(onConnect _iserver.FnOnConnect) {
	self.eventHandler.fnOnConnect = onConnect
}

func (self *WSServer) OnDisconnect(onDisconnect _iserver.FnOnDisconnect) {
	self.eventHandler.fnOnDisconnect = onDisconnect
}

func (self *WSServer) OnPing(onPing _iserver.FnOnPing) {
	self.eventHandler.fnOnPing = onPing
}

func (self *WSServer) OnMessage(onMessage _iserver.FnOnMessage) {
	self.eventHandler.fnOnMessage = onMessage
}