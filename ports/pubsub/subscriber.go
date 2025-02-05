package _ipubsub

type Subscriber[T any] interface {
	Receive() <-chan T
	Close()
}
