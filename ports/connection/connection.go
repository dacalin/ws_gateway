package _iconnection

import _connection_id "github.com/dacalin/ws_gateway/models/connection_id"

type Connection interface {
	Send(data []byte) // received a text/binary frame
	ConnectionId() _connection_id.ConnectionId
}
