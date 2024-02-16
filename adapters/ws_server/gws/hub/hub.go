package _gws_hub

import (
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iconnection "github.com/dacalin/ws_gateway/ports/connection"
	_ihub "github.com/dacalin/ws_gateway/ports/hub"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"sync"
)

var _ _ihub.Hub = (*Hub)(nil)

type Hub struct {
	connections map[_connection_id.ConnectionId]ConnectionData
	pubsub      _ipubsub.Client
}

var lock = &sync.Mutex{}
var instance *Hub

func Instance() *Hub {
	return instance
}

func New(pubsub _ipubsub.Client) *Hub {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = &Hub{
			connections: make(map[_connection_id.ConnectionId]ConnectionData),
			pubsub:      pubsub,
		}
	}

	return instance
}

func listener(data ConnectionData, pubsub _ipubsub.Client) {
	cid := data.connection.ConnectionId()
	subscriber := pubsub.Subscribe(cid.Value())

	for {
		select {
		case signal := <-data.endSignal:
			if signal == true {
				subscriber.Close()
				return
			}

		case msg := <-subscriber.Receive():
			data.connection.Send(msg)

		}
	}

}

func (self *Hub) Set(cid _connection_id.ConnectionId, conn _iconnection.Connection) {

	channel := make(chan bool)

	data := ConnectionData{
		endSignal:  channel,
		connection: conn,
	}

	self.connections[cid] = data

	go listener(data, self.pubsub)
}

func (self *Hub) Get(cid _connection_id.ConnectionId) (_iconnection.Connection, bool) {
	conn, found := self.connections[cid]

	if found == false {
		return nil, found
	}

	return conn.connection, found
}

func (self *Hub) Delete(cid _connection_id.ConnectionId) {
	conn, found := self.connections[cid]

	if found == false {
		return
	}

	conn.endSignal <- true

	delete(self.connections, cid)
}

func (self *Hub) PubSub() _ipubsub.Client {
	return self.pubsub
}

func (self *Hub) Send(cid _connection_id.ConnectionId, data []byte) {
	conn, found := self.Get(cid)

	if found == false {
		self.PubSub().Publish(cid.Value(), data)
	} else {
		conn.Send(data)
	}
}
