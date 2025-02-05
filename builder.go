package ws_gateway

import (
	"context"
	"errors"
	_pubsub "github.com/dacalin/ws_gateway/adapters/pubsub/redis"
	_gws_lib "github.com/dacalin/ws_gateway/adapters/ws_server/gws"
	_gws_hub "github.com/dacalin/ws_gateway/adapters/ws_server/gws/hub"
	"github.com/dacalin/ws_gateway/gateway"
	_logger "github.com/dacalin/ws_gateway/logger"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	_ipubsub "github.com/dacalin/ws_gateway/ports/pubsub"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"strings"
	"time"
)

func configPubSubDriver(config Config, ctx context.Context) (_ipubsub.Client[*redis.Message], error) {
	redisAddress := config.GWSDriver.PubSub.Host + ":" + strconv.Itoa(config.GWSDriver.PubSub.Port)

	var redisClient = redis.NewClient(&redis.Options{
		Addr:        redisAddress,
		Password:    config.GWSDriver.PubSub.Password,
		Username:    config.GWSDriver.PubSub.User,
		ReadTimeout: 0,
		PoolSize:    100,
		PoolTimeout: 60 * time.Second,
	})

	cmdResp := redisClient.Ping(ctx)
	if cmdResp.Err() != nil {
		log.Fatal(cmdResp.Err())
		return nil, cmdResp.Err()
	}

	pubsubClient := _pubsub.NewClient(redisClient, ctx)

	return pubsubClient, nil

}

func configGWSDriver(config Config, ctx context.Context) (_iserver.Server, _igateway.Gateway, error) {
	_logger.New(config.EnableDebugLog)

	var pubsubClient, pubsubErr = configPubSubDriver(config, ctx)
	if pubsubErr != nil {
		return nil, nil, pubsubErr
	}

	hub := _gws_hub.New(pubsubClient)
	connectionGateway := gateway.New(hub)

	server := _gws_lib.Create(config.GWSDriver.WSRoute, config.GWSDriver.PingIntervalSeconds, pubsubClient, config.CertFile, config.KeyFile)

	return server, connectionGateway, nil
}

func Create(config Config, ctx context.Context) (_iserver.Server, _igateway.Gateway, error) {

	switch strings.ToUpper(config.Driver) {

	case "GWS":
		server, connGW, err := configGWSDriver(config, ctx)
		return server, connGW, err

	default:
		return nil, nil, errors.New("WSGateway::Unsupported Driver " + config.Driver)

	}
}
