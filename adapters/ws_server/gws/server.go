package _gws_lib

import (
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	"crypto/tls"
	"fmt"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/lxzan/gws"
	"log"
	"net/http"
	"strconv"
	"time"
)

// WSServer represents a WebSocket server with generic type T (the pubsub message type), implementing the IServer interface.
type WSServer[T any] struct {
	_iserver.Server
	connectionRoute string
	eventHandler    EventHandler
	pubsub          _ipubsub.Client[T]
	certFile        string
	keyFile         string
}

// Create initializes a new WebSocket server with the given parameters.
func Create[T any](connectionRoute string, pingInterval int, pubsub _ipubsub.Client[T], hub _iserver.Hub, certFile string, keyFile string) *WSServer[T] {
	duration := time.Duration(pingInterval) * time.Second
	log.Println(fmt.Sprintf("Ping interval will close after %d seconds of inactivity", pingInterval))

	eventHandler := EventHandler{
		hub:            hub,
		pingInterval:   duration,
		fnOnConnect:    nil,
		fnOnDisconnect: nil,
		fnOnPing:       nil,
		fnOnMessage:    nil,
	}

	return &WSServer[T]{
		eventHandler:    eventHandler,
		connectionRoute: connectionRoute,
		pubsub:          pubsub,
		certFile:        certFile,
		keyFile:         keyFile,
	}
}

// Run starts the WebSocket server on the given port.
func (s *WSServer[T]) Run(port int) {
	upgrader := gws.NewUpgrader(&s.eventHandler, &gws.ServerOption{
		ReadAsyncEnabled: true,         // Parallel messages processing
		CompressEnabled:  true,         // Enable compression
		Recovery:         gws.Recovery, // Exception recovery
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/"+s.connectionRoute, func(writer http.ResponseWriter, request *http.Request) {
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

	if s.certFile != "" && s.keyFile != "" {
		// Start HTTPS server with TLS
		log.Println("Starting secure WebSocket server on wss://0.0.0.0" + addr)
		server := &http.Server{
			Addr:      addr,
			Handler:   mux,
			TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		}
		log.Fatal(server.ListenAndServeTLS(s.certFile, s.keyFile))
	} else {
		// Start HTTP server
		log.Println("Starting WebSocket server on ws://0.0.0.0" + addr)
		log.Fatal(http.ListenAndServe(addr, mux))
	}
}

// OnConnect sets the function to be called when a new connection is established.
func (s *WSServer[T]) OnConnect(onConnect _iserver.FnOnConnect) {
	s.eventHandler.fnOnConnect = onConnect
}

// OnDisconnect sets the function to be called when a connection is closed.
func (s *WSServer[T]) OnDisconnect(onDisconnect _iserver.FnOnDisconnect) {
	s.eventHandler.fnOnDisconnect = onDisconnect
}

// OnPing sets the function to be called when a ping is received.
func (s *WSServer[T]) OnPing(onPing _iserver.FnOnPing) {
	s.eventHandler.fnOnPing = onPing
}

// OnMessage sets the function to be called when a message is received.
func (s *WSServer[T]) OnMessage(onMessage _iserver.FnOnMessage) {
	s.eventHandler.fnOnMessage = onMessage
}
