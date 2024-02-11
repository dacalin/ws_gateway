package _igateway

import _connection_id "github.com/dacalin/ws_gateway/models/connection_id"

type Gateway interface {
	Send(id _connection_id.ConnectionId, data []byte)
	Broadcast(group string, data []byte)
	SetGroup(id _connection_id.ConnectionId, group string)
	RemoveGroup(id _connection_id.ConnectionId, group string)
}
