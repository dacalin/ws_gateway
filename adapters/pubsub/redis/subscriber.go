package _pubsub

import (
	_logger "github.com/dacalin/ws_gateway/logger"
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

func (s *Subscriber) Receive() chan []byte {
    chOut := make(chan []byte, 100)
    pubChan := s.subscriber.Channel() // returns a <-chan *redis.Message

    go func() {
        defer close(chOut)

        for {
            select {
            case <-s.endSignal:
                _logger.Instance().Println("Received end signal, stopping subscriber loop.")
                return

            case msg, ok := <-pubChan:
                if !ok {
                    // This means the subscription channel got closed
                    _logger.Instance().Println("Subscriber channel closed, exiting loop.")
                    return
                }
                // Process Redis message
                chOut <- []byte(msg.Payload)
                _logger.Instance().Println("New PubSub MSG:", msg.Payload)
            }
        }
    }()

    return chOut
}

func (self *Subscriber) Close() {
	_logger.Instance().Println("Subscriber Close.")
	self.endSignal <- true
	self.subscriber.Close()
}
