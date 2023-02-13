# GO STREAMING
- Recieve data from Redis Pub/sub, then push data to client.
- Should separate redis to improve performance
- Support SSE and WEBSOCKET

Library:
- gin: https://github.com/gin-gonic/gin
- redis pub/sub: https://github.com/go-redis/redis/v8
- websocket: https://github.com/gorilla/websocket
- broadcast: https://github.com/dustin/go-broadcast

# STRUCTOR
```
├── Dockerfile
├── README.MD
├── build.sh
├── config
│   └── streaming-api.json
├── go.mod
├── go.sum
├── main.go
├── model
│   ├── channels.go
│   └── config.go
└── utils
    └── helper.go
```

# MODEL
[PROCESSING] -> [REDIS PUB/SUB] -> [STREAMING] -> [CLIENT]

# CONFIG
- port: port streaming service listener
- redis: config redis pubsub
    - host: Host redis
    - port: Port redis
    - db: Database
    - password: Password

```
{
    "port": "8080",
    "redis": {
        "host": "localhost",
        "port": "6379",
        "db": "15",
        "password": "password"
    }
}
```

# INFO
- / : Get service info
- /admin : Get sse service info

# SSE STREAMING
- /streaming/:path/:channel : get sse streaming by path and channel
- Example:
    - /streaming/price/AAPL,IBM -> recieve data from redis pub/sub: price:AAPL and price:IBM
    - /streaming/account/AAPL,IBM -> recieve data from redis pub/sub: account:AAPL and account:IBM

# WS STREAMING
- /ws-streaming/:path/:channel : get sse streaming by path and channel
- Example:
    - /ws-streaming/price/AAPL,IBM -> recieve data from redis pub/sub: price:AAPL and price:IBM
    - /ws-streaming/account/AAPL,IBM -> recieve data from redis pub/sub: account:AAPL  and account:IBM


# TESTING
- Install package wscat
```
npm install -g wscat
```
- Run Websocket on terminal
```
wscat -c wss://localhost:8080/ws-streaming/price/AAPL,IBM
```

- Run Curl SSE on terminal 
```
curl http://localhost:8080/streaming/price/AAPL,IBM
```