# WS gateway: Simple, fast and reliable horizontal scaling of websockets.

## Introduction
WS Gateway is a simple, fast and reliable solution that allow you to **scale websockets** **server horizontally** easily, and painlessly. The plan is to offer different providers for the WS server and pubsub service, but right now the only option is GWS (https://github.com/lxzan/gws) as the WS Server and Redis as pubsub provider. 

A common solution to scale websockets is show is this diagram:

![image](https://github.com/dacalin/ws_gateway/assets/18492281/16d2ad56-173e-4998-86f7-93401a35a92e)

However, this library provide an abstraction that allow you to code like if you only had one instance. (See the example below)
![image](https://github.com/dacalin/ws_gateway/assets/18492281/027a6258-5064-47bb-9533-8a7681a59655)


## Public Interfaces

### Server
```
type FnOnConnect func(connectionId _connection_id.ConnectionId, params map[string]string)
type FnOnDisconnect func(connectionId _connection_id.ConnectionId)
type FnOnPing func(connectionId _connection_id.ConnectionId)
type FnOnMessage func(connectionId _connection_id.ConnectionId, data []byte)

type Server interface {
	Run(port int)
	OnConnect(FnOnConnect)
	OnDisconnect(FnOnDisconnect)
	OnPing(FnOnPing)
	OnMessage(FnOnMessage)
}
```


### Gateway
```
type Gateway interface {
	Send(id _connection_id.ConnectionId, data []byte)
	Broadcast(group string, data []byte)
	SetGroup(id _connection_id.ConnectionId, group string)
	RemoveGroup(id _connection_id.ConnectionId, group string)
}
```


## Example
An example app can be found here https://github.com/dacalin/demo_chat

```
	WSConfig := ws_gateway.Config{
		Driver:         ws_gateway.DRIVER_WS_GWS,
		EnableDebugLog: true,
		GWSDriver: ws_gateway.GWSDriverConfig{
			PubSub: 	ws_gateway.PubSubDriverConfig{
				Driver: 		ws_gateway.DRIVER_PUBSUB_REDIS,
				Host:           config.RedisHost,
				Port:           config.RedisPort,
			},
			PingIntervalSeconds: config.WsPingIntervalSeconds,
			WSRoute:             "connect",
		},
	}

	wsServer, wsGatewayConnectio1, err := ws_gateway.Create(WSConfig, ctx)
	if err != nil {
		panic(err)
	}

	wsServer.OnConnect(func(connectionId _connection_id.ConnectionId, params map[string]string) {
		wsGatewayConnection.SetGroup(connectionId, "demo-room")
	})

  // On a message, broadcast the message to all clients. This will automatically
  // sync with different instances through the pubsub service.
	wsServer.OnMessage(
		func(connectionId _connection_id.ConnectionId, data []byte) {
			wsGatewayConnection.Broadcast("demo-room", data)
		})


 // Run the server
	wsServer.Run(config.WsPort)

```

## Roadmap
The lib was designed thinking in future support of different services. 
The next adapters to be implemented are:
- AWS API Gateway
