package _pubsub

import (
	"context"
	_logger "github.com/dacalin/ws_gateway/logger"
	_ipubsub "github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
	"log"
)

var _ _ipubsub.Client[*redis.Message] = (*Client)(nil)

type Client struct {
	_ipubsub.Client[*redis.Message]
	client *redis.Client
	ctx    context.Context
}

func NewClient(client *redis.Client, ctx context.Context) _ipubsub.Client[*redis.Message] {
	return &Client{
		client: client,
		ctx:    ctx,
	}
}

func (self *Client) Subscribe(topics ...string) _ipubsub.Subscriber[*redis.Message] {
	subscriber := self.client.Subscribe(self.ctx, topics...)
	return NewSubscriber(subscriber, self.ctx)
}

func (self *Client) Publish(topic string, message []byte) {
	_logger.Instance().Printf("Publish, topic=%s, msg=%s\n", topic, string(message))

	cmd := self.client.Publish(self.ctx, topic, message)
	if cmd != nil && cmd.Err() != nil {
		log.Fatal(cmd.Err())
	} else {
		num, _ := cmd.Result()

		_logger.Instance().Printf("Publish, listeners=%d", num)
	}

}
