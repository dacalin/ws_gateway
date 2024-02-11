package _pubsub

import (
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
)
import "context"

type Subscriber struct {
	_ipubsub.Subscriber
	client     *redis.Client
	ctx        context.Context
	subscriber *redis.PubSub
}

func NewSubscriber(subscriber *redis.PubSub, ctx context.Context) *Subscriber {
	return &Subscriber{
		subscriber: subscriber,
		ctx:        ctx,
	}
}

func (self *Subscriber) Receive() ([]byte, error) {
	msg, err := self.subscriber.ReceiveMessage(self.ctx)
	if err != nil {
		return nil, err
	}

	return []byte(msg.Payload), nil
}

func (self *Subscriber) Close() {
	self.subscriber.Close()
}
