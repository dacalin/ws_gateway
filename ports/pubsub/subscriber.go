package _ipubsub

type Subscriber interface {
	Receive() ([]byte, error)
	Close()
}
