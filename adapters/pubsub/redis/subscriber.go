package _pubsub

import (
	_logger "github.com/dacalin/ws_gateway/logger"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
)
import "context"

var _ _ipubsub.Subscriber[*redis.Message] = (*Subscriber)(nil)

type Subscriber struct {
	_ipubsub.Subscriber[*redis.Message]
	ctx        context.Context
	subscriber *redis.PubSub
}

func NewSubscriber(subscriber *redis.PubSub, ctx context.Context) *Subscriber {
	return &Subscriber{
		subscriber: subscriber,
		ctx:        ctx,
	}
}

func (s *Subscriber) Receive() <-chan *redis.Message {
	return s.subscriber.Channel()
}

func (self *Subscriber) Close() {
	_logger.Instance().Println("Subscriber Close.")
	self.subscriber.Close()
}
