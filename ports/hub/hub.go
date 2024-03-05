package _ihub

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_client_connection "github.com/dacalin/ws_gateway/ports/connection"
)

type Hub interface {
	Set(cid _connection_id.ConnectionId, conn _client_connection.Connection)
	Get(cid _connection_id.ConnectionId) (_client_connection.Connection, bool)
	Delete(cid _connection_id.ConnectionId)
	Send(cid _connection_id.ConnectionId, data []byte)
	SendTo(topic string, data []byte)
	ListenTo(cid _connection_id.ConnectionId, topic string)
}
