package _gws_hub

import (
	_client_connection "github.com/dacalin/ws_gateway/ports/connection"
)

type ConnectionData struct {
	endSignal  chan bool
	connection _client_connection.Connection
}
