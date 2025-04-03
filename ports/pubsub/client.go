package _ipubsub

type Client[T any] interface {
	Subscribe(channels ...string) Subscriber[T]
	Publish(channel string, message []byte)
	GetNumSub(channel string) int64
}
