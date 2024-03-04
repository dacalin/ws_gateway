package _pubsub

import (
	_logger "github.com/dacalin/ws_gateway/logger"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	"github.com/go-redis/redis/v8"
	"log"
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
	chOut := make(chan []byte, 100)

	go func() {

		for {
			msgi, err := self.subscriber.Receive(self.ctx)
			if err != nil {
				log.Fatal("Received Redis Error. ", err.Error())
			} else {
				switch msg := msgi.(type) {
				case *redis.Message:
					chOut <- []byte(msg.Payload)
					_logger.Instance().Println("New PubSub MSG")
					<-chOut

				default:
					_logger.Instance().Println("New PubSub Control MSG")
				}
			}

		}
	}()

	return chOut

}

func (self *Subscriber) Close() {
	self.endSignal <- true
	self.subscriber.Close()
}
