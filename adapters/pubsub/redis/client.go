package _pubsub

import (
	"context"
	_ipubsub "github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
	"log"
)

var _ _ipubsub.Client = (*Client)(nil)

type Client struct {
	_ipubsub.Client
	client *redis.Client
	ctx    context.Context
}

func NewClient(client *redis.Client, ctx context.Context) _ipubsub.Client {
	return &Client{
		client: client,
		ctx:    ctx,
	}
}

func (self *Client) Subscribe(channels ...string) _ipubsub.Subscriber {
	subscriber := self.client.Subscribe(self.ctx, channels...)
	return NewSubscriber(subscriber, self.ctx)
}

func (self *Client) Publish(channel string, message []byte) {
	cmd := self.client.Publish(self.ctx, channel, message)
	if cmd != nil && cmd.Err() != nil {
		log.Fatal(cmd.Err())
	}
}
