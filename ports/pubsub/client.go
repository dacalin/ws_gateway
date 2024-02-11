package _ipubsub

type Client interface {
	Subscribe(channels ...string) Subscriber
	Publish(channel string, message []byte)
}
