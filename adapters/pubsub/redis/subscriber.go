package _pubsub

import (
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
)
import "context"

var _ _ipubsub.Subscriber = (*Subscriber)(nil)

type Subscriber struct {
	_ipubsub.Subscriber
	client     *redis.Client
	ctx        context.Context
	subscriber *redis.PubSub
	endSignal  chan bool
}

func NewSubscriber(subscriber *redis.PubSub, ctx context.Context) *Subscriber {
	endSignal := make(chan bool)
	return &Subscriber{
		subscriber: subscriber,
		ctx:        ctx,
		endSignal:  endSignal,
	}
}

func (self *Subscriber) Receive() chan []byte {
	ch := self.subscriber.Channel()

	chOut := make(chan []byte)

	go func() {
		// go routine to wrap the redis type and create a different channel type that match the interface
		for {
			select {
			case msg := <-ch:
				chOut <- []byte(msg.Payload)

			case <-self.endSignal:
				println("Receive go routine end")
				return

			}
		}
	}()

	return chOut

}

func (self *Subscriber) Close() {
	self.endSignal <- true
	self.subscriber.Close()
}
