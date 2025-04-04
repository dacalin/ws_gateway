package _pubsub

import (
	"context"
	evbus "github.com/asaskevich/EventBus"
	_logger "github.com/dacalin/ws_gateway/logger"
	_ipubsub "github.com/dacalin/ws_gateway/ports/pubsub"
)

var _ _ipubsub.Client[*[]byte] = (*Client)(nil)

// Client represents a pubsub client.
type Client struct {
	_ipubsub.Client[*[]byte]
	ctx context.Context
	bus *evbus.Bus
}

// NewClient creates a new pubsub client with the given context.
func NewClient(ctx context.Context, bus *evbus.Bus) _ipubsub.Client[*[]byte] {
	return &Client{
		ctx: ctx,
		bus: bus}
}

// Subscribe subscribes to the given topics.
func (c *Client) Subscribe(topics ...string) _ipubsub.Subscriber[*[]byte] {
	return NewSubscriber(c.ctx, c.bus, topics...)

}

// Publish publishes a message to the given topic.
func (c *Client) Publish(topic string, message []byte) {
	_logger.Instance().Printf("Publish, topic=%s, msg=%s\n", topic, string(message))
	(*c.bus).Publish(topic, message)
}

func (c *Client) IsListened(topic string) bool {
	return (*c.bus).HasCallback(topic)
}
