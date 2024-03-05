package _gws_hub

import (
	"context"
	_client_connection "github.com/dacalin/ws_gateway/ports/connection"
)

type ConnectionData struct {
	ctx        context.Context
	cancel     context.CancelFunc
	connection _client_connection.Connection
}
