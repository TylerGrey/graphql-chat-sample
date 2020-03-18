# Simple Graphql chatting server with Redis

## 1. Requirements
* Docker
* `go env -w GO111MODULE=auto`
* go-bindata

## 2. Run
```shell script
make redis
make schema
make run
```

http://localhost:8080/
```graphql
{
  ping
}
```

## 3. Chat
### Subscriber
```graphql
subscription {
  onMessage {
    id
    msg
    createdAt
  }
}
```

### Publisher
```graphql
mutation {
  sendMessage(msg: "Hello, World!") {
    id
    msg
    createdAt
  }
}
```