package ws_gateway

import (
	"context"
	_pubsub "github.com/dacalin/ws_gateway/adapters/pubsub/redis"
	_gws_lib "github.com/dacalin/ws_gateway/adapters/ws_server/gws"
	"github.com/dacalin/ws_gateway/gateway"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	"github.com/go-redis/redis/v8"
)

func CreateServer(redisAddress string, pingIntervalSeconds int, ctx context.Context) _iserver.Server {
	var redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	pubsubClient := _pubsub.NewClient(redisClient, ctx)

	server := _gws_lib.Create("connect", pingIntervalSeconds, pubsubClient)

	return &server
}

func CreateConnectionGateway() _igateway.Gateway {
	return gateway.New()
}
