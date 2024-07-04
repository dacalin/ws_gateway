package _gws_lib

import (
	"crypto/tls"
	"fmt"
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
	certFile        string
	keyFile         string
}

func Create(connectionRoute string, pingInterval int, pubsub _ipubsub.Client, certFile string, keyFile string) *WSServer {
	duration := time.Duration(pingInterval) * time.Second
	log.Println(fmt.Sprintf("Ping interval will close after %d seconds of inactivity", pingInterval))

	eventHandler := EventHandler{
		pingInterval:   duration,
		fnOnConnect:    nil,
		fnOnDisconnect: nil,
		fnOnPing:       nil,
		fnOnMessage:    nil,
	}

	return &WSServer{
		eventHandler:    eventHandler,
		connectionRoute: connectionRoute,
		pubsub:          pubsub,
		certFile:        certFile,
		keyFile:         keyFile,
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

	addr := ":" + strconv.Itoa(port)

	if self.certFile != "" && self.keyFile != "" {
		// Start HTTPS server with TLS
		log.Println("Starting secure WebSocket server on wss://0.0.0.0" + addr)
		server := &http.Server{
			Addr:      addr,
			Handler:   mux,
			TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		}
		log.Fatal(server.ListenAndServeTLS(self.certFile, self.keyFile))
	} else {
		// Start HTTP server
		log.Println("Starting WebSocket server on ws://0.0.0.0" + addr)
		log.Fatal(http.ListenAndServe(addr, mux))
	}
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
