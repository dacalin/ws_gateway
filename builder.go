package ws_gateway

import (
	"context"
	"errors"
	_pubsub "github.com/dacalin/ws_gateway/adapters/pubsub/redis"
	_gws_lib "github.com/dacalin/ws_gateway/adapters/ws_server/gws"
	_gws_hub "github.com/dacalin/ws_gateway/adapters/ws_server/gws/hub"
	"github.com/dacalin/ws_gateway/gateway"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	_ipubsub "github.com/dacalin/ws_gateway/ports/pubsub"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

func configPubSubDriver(config Config, ctx context.Context) _ipubsub.Client {
	redisAddress := config.GWSDriver.PubSub.Host + ":" + strconv.Itoa(config.GWSDriver.PubSub.Port)

	var redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	pubsubClient := _pubsub.NewClient(redisClient, ctx)

	return pubsubClient

}

func configGWSDriver(config Config, ctx context.Context) (_iserver.Server, _igateway.Gateway) {
	var pubsubClient = configPubSubDriver(config, ctx)

	hub := _gws_hub.New(pubsubClient)
	connectionGateway := gateway.New(hub)

	server := _gws_lib.Create(config.GWSDriver.WSRoute, config.GWSDriver.PingIntervalSeconds, pubsubClient, config.EnableDebugLog)
	return server, connectionGateway
}

func Create(config Config, ctx context.Context) (_iserver.Server, _igateway.Gateway, error) {

	switch strings.ToUpper(config.Driver) {

	case "GWS":
		server, connGW := configGWSDriver(config, ctx)
		return server, connGW, nil

	default:
		return nil, nil, errors.New("WSGateway::Unsupported Driver " + config.Driver)

	}
}
