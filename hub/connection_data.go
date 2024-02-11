package _hub

import (
	_client_connection "github.com/dacalin/ws_gateway/ports/connection"
)

type ConnectionData struct {
	channel    chan string
	connection _client_connection.Connection
}
