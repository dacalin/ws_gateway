package _gws_hub

import (
	"context"
	"encoding/json"
	"sync"
	_logger "github.com/dacalin/ws_gateway/logger"
	_connection_id "github.com/dacalin/ws_gateway/models/connection_id"
	_iconnection "github.com/dacalin/ws_gateway/ports/connection"
	"github.com/dacalin/ws_gateway/ports/pubsub"
)

// Hub keep track of connections, subscriptions and pubsub clients.
type Hub[T any] struct {
	connections sync.Map
	pubsub _ipubsub.Client[T]
}

// Returns a new Hub instance.
func New[T any] (pubsub _ipubsub.Client[T]) *Hub[T] {
	return &Hub[T]{pubsub: pubsub}
}

// Converts the given message to a byte slice.
func convert[T any](msg T) []byte {
	_logger.Instance().Printf("convert")

	var data []byte
	switch v := any(msg).(type) {
	case *[]byte:		
	    _logger.Instance().Printf("*[]byte",)
 
		if v == nil {
			_logger.Instance().Printf("*[]byte is nil")
			data = nil
		} else {
			data = []byte(string(*v))
			_logger.Instance().Printf("data=%v", data)
		}

	case []byte:
		// If it's already a []byte, use it as is.
		data = v
		_logger.Instance().Printf("data=", data)
	case string:
		// Convert a string to []byte.
		data = []byte(v)
		_logger.Instance().Printf("datas=", data)

	default:
		_logger.Instance().Printf("default")

		// For other types, attempt to marshal the value to JSON.
		var err error
		data, err = json.Marshal(v)
		if err != nil {
			_logger.Instance().Printf("Failed to marshal message: %v", err)
			return data
		}
	}

	return data
}

// listener listens for messages on the given topic and sends them to the
func (h *Hub[T]) listener(data ConnectionData, pubsub _ipubsub.Client[T], topic string) {
	_logger.Instance().Printf("listening cid=%s topic=%s", data.connection.ConnectionId(), topic)
	defer _logger.Instance().Printf("quit listener end cid=%s topic=%s", data.connection.ConnectionId(), topic)

	subscriber := pubsub.Subscribe(topic)

	for {
		select {
		case <-data.ctx.Done():
			_logger.Instance().Printf("listener end signal cid=%s, topic=%s", data.connection.ConnectionId(), topic)
			subscriber.Close()
			return

		case msg := <-subscriber.Receive():
			msgT := convert(msg)
			_logger.Instance().Printf("RECEIVER MSG cid=%s, MSG=%s", data.connection.ConnectionId(), string(msgT))
			go data.connection.Send(msgT)

		}
	}
}

// Set adds a connection to the hub.
func (h  *Hub[T]) Set(cid _connection_id.ConnectionId, conn _iconnection.Connection) {

	ctx, cancel := context.WithCancel(context.Background())

	data := ConnectionData{
		ctx:        ctx,
		cancel:     cancel,
		connection: conn,
	}

	h.connections.Store(cid, data)

	go h.listener(data, h.pubsub, cid.Value())
}

// Get returns the connection with the given connection ID.
func (h  *Hub[T]) Get(cid _connection_id.ConnectionId) (_iconnection.Connection, bool) {
	connDataI, found := h.connections.Load(cid)

	if found == false {
		return nil, found
	}

	connData := connDataI.(ConnectionData)
	return connData.connection, found
}

// Delete removes the connection with the given connection ID.
func (h  *Hub[T]) Delete(cid _connection_id.ConnectionId) {
	connDataI, found := h.connections.Load(cid)
	if found == false {
		_logger.Instance().Println("Hub Delete: cid not found")
		return
	}

	connData := connDataI.(ConnectionData)
	connData.cancel()

	h.connections.Delete(cid)
}

// PubSub returns the pubsub client.
func (h  *Hub[T]) PubSub() _ipubsub.Client[T] {
	return h.pubsub
}

// Sends the given data to the connection with the given connection ID.
func (h  *Hub[T]) Send(cid _connection_id.ConnectionId, data []byte) {
	_logger.Instance().Printf("Send To cid=%s", cid.Value())

	conn, found := h.Get(cid)

	if found == false {
		_logger.Instance().Println("Send using PubSub")
		h.PubSub().Publish(cid.Value(), data)
	} else {
		conn.Send(data)
	}
}

// SendTo sends the given data to the given topic.
func (h  *Hub[T]) SendTo(topic string, data []byte) {
	_logger.Instance().Printf("Send To topic=%s, msg=%s", topic, data)
	h.PubSub().Publish(topic, data)
}

// ListenTo start a listener for the given connection ID and topic.
func (h  *Hub[T]) ListenTo(cid _connection_id.ConnectionId, topic string) {
	connDataI, found := h.connections.Load(cid)
	connData := connDataI.(ConnectionData)

	if found == true {
		_logger.Instance().Printf("Listen To topic=%s, cid%s", topic, cid)
		go h.listener(connData, h.pubsub, topic)
	}
}
