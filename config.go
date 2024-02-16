package ws_gateway

//redisAddress string, pingIntervalSeconds int,

type Config struct {
	Driver         string
	GWSDriver      GWSDriverConfig
	EnableDebugLog bool
}

type GWSDriverConfig struct {
	RedisHost           string
	RedisPort           int
	PingIntervalSeconds int
	WSRoute             string
}
