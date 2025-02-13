package ws_gateway

import (
	"log"
	"strconv"
	"strings"
	"time"
	"context"
	"errors"
	_redis_pubsub "github.com/dacalin/ws_gateway/adapters/pubsub/redis"
	_inprocess_pubsub "github.com/dacalin/ws_gateway/adapters/pubsub/inprocess"
	_gws_lib "github.com/dacalin/ws_gateway/adapters/ws_server/gws"
	_gws_hub "github.com/dacalin/ws_gateway/adapters/ws_server/gws/hub"
	_logger "github.com/dacalin/ws_gateway/logger"
	_igateway "github.com/dacalin/ws_gateway/ports/gateway"
	_iserver "github.com/dacalin/ws_gateway/ports/server"
	evbus "github.com/asaskevich/EventBus"
	"github.com/dacalin/ws_gateway/gateway"
	"github.com/go-redis/redis/v8"
)

// Config websocket server with the corresponding pubsub driver.
func configGWSDriver(ctx context.Context, config Config) (_iserver.Server, _igateway.Gateway, error) {
	_logger.New(config.EnableDebugLog)

	var server _iserver.Server
	var connectionGateway _igateway.Gateway

	switch strings.ToUpper(config.GWSDriver.PubSub.Driver) {
		case DRIVER_PUBSUB_REDIS:
			// Connect to Redis
			redisAddress := config.GWSDriver.PubSub.Host + ":" + strconv.Itoa(config.GWSDriver.PubSub.Port)
			log.Println("Configuring PubSub::Connecting to Redis at " + redisAddress)
			
			var redisClient = redis.NewClient(&redis.Options{
				Addr:        redisAddress,
				Password:    config.GWSDriver.PubSub.Password,
				Username:    config.GWSDriver.PubSub.User,
				ReadTimeout: 0,
				PoolSize:    100,
				PoolTimeout: 60 * time.Second,
			})

			// Ping Redis, check if connection is established
			cmdResp := redisClient.Ping(ctx)
			if cmdResp.Err() != nil {
				log.Fatal(cmdResp.Err())
				return nil, nil, cmdResp.Err()
			}

			// Create PubSub client
			pubsubClient := _redis_pubsub.NewClient(redisClient, ctx)

			// Create hub to manage connections
			hub := _gws_hub.New(pubsubClient)

			// Create connection gateway
			connectionGateway = gateway.New(hub)

			// Create server
			server = _gws_lib.Create(config.GWSDriver.WSRoute, config.GWSDriver.PingIntervalSeconds, pubsubClient, hub, config.CertFile, config.KeyFile)


		 default:
			log.Println("Configuring PubSub::Setting up in-process PubSub")
			// Create in-process PubSub client
			bus := evbus.New()
			pubsubClient := _inprocess_pubsub.NewClient(ctx, &bus)

			// Create hub to manage connections
			hub := _gws_hub.New(pubsubClient)
			// Create connection gateway
			connectionGateway = gateway.New(hub)
			// Create server
			server = _gws_lib.Create(config.GWSDriver.WSRoute, config.GWSDriver.PingIntervalSeconds, pubsubClient, hub, config.CertFile, config.KeyFile)
	}


	return server, connectionGateway, nil
}

// Create creates a new instance of the gateway and server.
func Create(config Config, ctx context.Context) (_iserver.Server, _igateway.Gateway, error) {

	switch strings.ToUpper(config.Driver) {

		case "GWS":
			server, connGW, err := configGWSDriver(ctx, config)
			return server, connGW, err

		default:
			return nil, nil, errors.New("WSGateway::Unsupported Driver " + config.Driver)

	}
}
