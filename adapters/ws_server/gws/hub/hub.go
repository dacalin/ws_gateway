package _gws_hub

import (
	"context"
	_logger "github.com/dacalin/ws_gateway/logger"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iconnection "github.com/dacalin/ws_gateway/ports/connection"
	_ihub "github.com/dacalin/ws_gateway/ports/hub"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
	"sync"
)

var _ _ihub.Hub = (*Hub)(nil)

type Hub struct {
	connections sync.Map
	//connections map[_connection_id.ConnectionId]ConnectionData
	pubsub _ipubsub.Client[*redis.Message]
}

var lock = &sync.Mutex{}
var instance *Hub

func Instance() *Hub {
	return instance
}

func New(pubsub _ipubsub.Client[*redis.Message]) *Hub {
	lock.Lock()
	defer lock.Unlock()
	if instance == nil {
		instance = &Hub{
			connections: sync.Map{},
			pubsub:      pubsub,
		}
	}

	return instance
}

func listener(data ConnectionData, pubsub _ipubsub.Client[*redis.Message], topic string) {
	_logger.Instance().Printf("listening cid=%s topic=%s", data.connection.ConnectionId(), topic)

	subscriber := pubsub.Subscribe(topic)

	for {
		select {
		case <-data.ctx.Done():
			_logger.Instance().Printf("listener end signal cid=%s", data.connection.ConnectionId())
			subscriber.Close()
			return

		case msg := <-subscriber.Receive():
			_logger.Instance().Printf("RECEIVER MSG cid=%s, MSG=%s", data.connection.ConnectionId(), msg)
			go data.connection.Send([]byte(msg.Payload))

		}
	}

}

func (self *Hub) Set(cid _connection_id.ConnectionId, conn _iconnection.Connection) {

	ctx, cancel := context.WithCancel(context.Background())

	data := ConnectionData{
		ctx:        ctx,
		cancel:     cancel,
		connection: conn,
	}

	self.connections.Store(cid, data)

	go listener(data, self.pubsub, cid.Value())
}

func (self *Hub) Get(cid _connection_id.ConnectionId) (_iconnection.Connection, bool) {
	connDataI, found := self.connections.Load(cid)

	if found == false {
		return nil, found
	}

	connData := connDataI.(ConnectionData)
	return connData.connection, found
}

func (self *Hub) Delete(cid _connection_id.ConnectionId) {
	connDataI, found := self.connections.Load(cid)
	if found == false {
		return
	}

	connData := connDataI.(ConnectionData)
	connData.cancel()

	self.connections.Delete(cid)
}

func (self *Hub) PubSub() _ipubsub.Client[*redis.Message] {
	return self.pubsub
}

func (self *Hub) Send(cid _connection_id.ConnectionId, data []byte) {
	_logger.Instance().Printf("Send To cid=%s", cid.Value())

	//conn, found := self.Get(cid)

	//if found == false {
	_logger.Instance().Println("Send using PubSub")
	self.PubSub().Publish(cid.Value(), data)
	//} else {
	//	conn.Send(data)
	//}
}

func (self *Hub) SendTo(topic string, data []byte) {
	_logger.Instance().Printf("Send To topic=%s", topic)
	self.PubSub().Publish(topic, data)
}

func (self *Hub) ListenTo(cid _connection_id.ConnectionId, topic string) {
	connDataI, found := self.connections.Load(cid)
	connData := connDataI.(ConnectionData)

	if found == true {
		_logger.Instance().Printf("Listen To topic=%s, cid%s", topic, cid)
		go listener(connData, self.pubsub, topic)
	}

}
