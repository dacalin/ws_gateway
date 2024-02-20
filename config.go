package ws_gateway

var DRIVER_WS_GWS = "gws"
var DRIVER_PUBSUB_REDIS = "redis"

type Config struct {
	Driver         string
	GWSDriver      GWSDriverConfig
	EnableDebugLog bool
}

type GWSDriverConfig struct {
	PubSub              PubSubDriverConfig
	PingIntervalSeconds int
	WSRoute             string
}

type PubSubDriverConfig struct {
	Driver string
	Host   string
	Port   int
}
