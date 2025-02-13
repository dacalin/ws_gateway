package _pubsub

import (
	"context"
	_logger "github.com/dacalin/ws_gateway/logger"
	"github.com/dacalin/ws_gateway/ports/pubsub"
	evbus "github.com/asaskevich/EventBus"

)

var _ _ipubsub.Subscriber[*[]byte] = (*Subscriber)(nil)

type Subscriber struct {
	_ipubsub.Subscriber[*[]byte]
	ctx        context.Context
	topics	  []string
	bus *evbus.Bus
	channel chan *[]byte
	busSubscriber []evbus.BusSubscriber
}

// NewSubscriber creates a new Subscriber with the given context and topics.
func NewSubscriber(ctx context.Context, bus *evbus.Bus, topics ...string) *Subscriber {
	_logger.Instance().Printf("NewSubscriber topics:%s", topics)

	var channel = make(chan *[]byte, 10)
	_logger.Instance().Printf("NewSubscriber::Channel=%s", channel)

	s := &Subscriber{
		ctx:        ctx,
		topics:     topics,
		bus: 		bus,
		channel: 	channel,
	}
	_logger.Instance().Printf("NewSubscriber::subscriber=%s", s)

	for _, topic := range topics {
		(*bus).Subscribe(topic, s.subscribe)
	}
	
	return s
}

func (s *Subscriber) subscribe(msg []byte){
	_logger.Instance().Printf("subscribe msg=%s\n", string(msg))
	_logger.Instance().Printf("subscribe::Channel=%s", s.channel)
	_logger.Instance().Printf("subscribe::Subscriber=%s", s)

	s.channel <- &msg
}

// Receive returns a channel to receive messages.
func (s *Subscriber) Receive() <-chan *[]byte {
	return s.channel
}

// Close shuts down the Subscriber.
func (s *Subscriber) Close() {
	_logger.Instance().Printf("Subscriber Close topics:%s", s.topics)

	for _, topic := range s.topics {
		(*s.bus).Unsubscribe(topic, s.subscribe)		
	}
	
	// close channel
	close(s.channel)
}
