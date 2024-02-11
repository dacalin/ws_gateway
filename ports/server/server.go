package _iserver

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
)

type FnOnConnect func(connectionId _connection_id.ConnectionId, params map[string]string)
type FnOnDisconnect func(connectionId _connection_id.ConnectionId)
type FnOnPing func(connectionId _connection_id.ConnectionId)
type FnOnMessage func(connectionId _connection_id.ConnectionId, data []byte)

type Server interface {
	Run(port int)
	OnConnect(FnOnConnect)
	OnDisconnect(FnOnDisconnect)
	OnPing(FnOnPing)
	OnMessage(FnOnMessage)
}
