# Game Example for dubbo-go

English | [中文](README_CN.md)

This example demonstrates a football game application using dubbo-go as an RPC framework. It consists of two services: a game service that handles game logic (login, scoring, ranking) and a gate service that provides HTTP endpoints for the web frontend and acts as a gateway between the frontend and game service.

## Architecture

- **game**: Game service that handles game logic (Login, Score, Rank)
- **gate**: Gate service that provides HTTP API for web frontend and RPC service for game service
- **website**: Web frontend for the football game
- **proto**: Protocol buffer definitions for game and gate services

## Contents

- `game/go-server/cmd/main.go` - Game service server implementation
- `game/pkg/provider.go` - Game service handler implementation
- `gate/go-server/cmd/main.go` - Gate service server with HTTP and RPC
- `gate/pkg/provider.go` - Gate service handler implementation
- `gate/pkg/consumer.go` - Game service client for gate service
- `proto/` - Protocol buffer definitions and generated code
- `website/` - Web frontend files

## How to run

### Run Game Server

Start the game service server (listens on port 20000):

```shell
go run ./game/go-server/cmd/main.go
```

### Run Gate Server

Start the gate service server (listens on port 20001 for RPC and 8089 for HTTP):

```shell
go run ./gate/go-server/cmd/main.go
```

### Access Web Frontend

Open `http://127.0.0.1:8089/` in your browser to access the game frontend. The frontend will communicate with the gate service HTTP API.

## Service Communication

1. **Frontend → Gate Service**: HTTP requests to `/login`, `/score`, `/rank`
2. **Gate Service → Game Service**: RPC calls using Triple protocol
3. **Game Service → Gate Service**: RPC calls for gate operations

## Testing

You can test the services using curl:

```shell
# Test login
curl "http://127.0.0.1:8089/login?name=player1"

# Test score
curl -X POST http://127.0.0.1:8089/score \
  -H "Content-Type: application/json" \
  -d '{"name":"player1","score":1}'

# Test rank
curl "http://127.0.0.1:8089/rank?name=player1"
```

## Notes

- Make sure both game server and gate server are running before accessing the web frontend
- The game server must be started before the gate server, as the gate server needs to connect to the game service
- Both services use Triple protocol for RPC communication
