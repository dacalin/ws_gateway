package _hub

import (
	"fmt"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_client_connection "github.com/dacalin/ws_gateway/ports/connection"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"sync"
)

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

func (self *Hub) listener(data ConnectionData) {
	cid := data.connection.ConnectionId()
	subscriber := self.pubsub.Subscribe(cid.Value())

	for {
		select {
		case msg := <-data.channel:
			if msg == "close" {
				println("end goroutine")
				return
			}

		default:
			msg, err := subscriber.Receive()
			if err != nil {
				fmt.Println(err)
			}

			conn, found := self.Get(cid)
			if found {
				conn.Send(msg)
			}
		}
	}
}

func (self *Hub) Set(cid _connection_id.ConnectionId, conn _client_connection.Connection) {

	channel := make(chan string)

	data := ConnectionData{
		channel:    channel,
		connection: conn,
	}

	self.connections[cid] = data

	go self.listener(data)
}

func (self *Hub) Get(cid _connection_id.ConnectionId) (_client_connection.Connection, bool) {
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

	conn.channel <- "close"

	delete(self.connections, cid)
}

func (self *Hub) PubSub() _ipubsub.Client {
	return self.pubsub
}
