package _ipubsub

type Client[T any] interface {
	Subscribe(channels ...string) Subscriber[T]
	Publish(channel string, message []byte)
	IsListened(channel string) bool
}
