package ws_gateway

var DRIVER_WS_GWS = "gws"
var DRIVER_PUBSUB_REDIS = "redis"
var DRIVER_PUBSUB_INTERNAL = "inprocess"

type Config struct {
	Driver         string
	GWSDriver      GWSDriverConfig
	EnableDebugLog bool
	CertFile       string
	KeyFile        string
}

type GWSDriverConfig struct {
	PubSub              PubSubDriverConfig
	PingIntervalSeconds int
	WSRoute             string
}

type PubSubDriverConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Driver   string
}
