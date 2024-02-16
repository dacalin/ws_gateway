package _ipubsub

type Subscriber interface {
	Receive() chan []byte
	Close()
}
