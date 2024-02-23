package _gws_lib

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	"github.com/lxzan/gws"
	"log"
	"net/http"
	"strconv"
	"time"
)

var _ _iserver.Server = (*WSServer)(nil)

type WSServer struct {
	_iserver.Server
	connectionRoute string
	eventHandler    EventHandler
	pubsub          _ipubsub.Client
}

func Create(connectionRoute string, pingInterval int, pubsub _ipubsub.Client, debug bool) *WSServer {
	duration := time.Duration(pingInterval) * time.Second

	eventHandler := EventHandler{
		pingInterval:   duration,
		fnOnConnect:    nil,
		fnOnDisconnect: nil,
		fnOnPing:       nil,
		fnOnMessage:    nil,
		debug:          debug,
	}

	return &WSServer{
		eventHandler:    eventHandler,
		connectionRoute: connectionRoute,
		pubsub:          pubsub,
	}
}

func (self *WSServer) Run(port int) {
	upgrader := gws.NewUpgrader(&self.eventHandler, &gws.ServerOption{
		ReadAsyncEnabled: true,         // Parallel messages processing
		CompressEnabled:  true,         // Enable compression
		Recovery:         gws.Recovery, // Exception recovery
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/"+self.connectionRoute, func(writer http.ResponseWriter, request *http.Request) {
		socket, err := upgrader.Upgrade(writer, request)
		if err != nil {
			return
		}

		cidParam := request.URL.Query().Get("cid")
		if cidParam == "" {
			log.Printf("GWSHandlerError::Expected cid param, got none.")
			socket.NetConn().Close()
			return
		}

		// Get connection Id
		cid := _connection_id.New(cidParam)
		socket.Session().Store("cid", cid)

		// Get query params
		queryParams := request.URL.Query()
		params := make(map[string]string)

		for key, value := range queryParams {
			params[key] = value[0]
		}

		socket.Session().Store("params", params)

		go func() {
			socket.ReadLoop() // Blocking prevents the context from being GC.
		}()
	})

	http.ListenAndServe(":"+strconv.Itoa(port), mux)
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
