# go-broker

Simple Message broker written in Go.

## Installation

Runnning the message broker is as easy as any other go app. Just use the command

```bash
go run cmd/go-broker/main.go
```

## APIs

### POST APIs

- Register Publisher: 

```bash
curl -X POST -d '{"id":"pub1"}' -H "Content-Type: application/json" http://localhost:8080/register-publisher

```

- Register Subscriber:

```bash
curl -X POST -d '{"id":"sub1", "topic":"topic1"}' -H "Content-Type: application/json" http://localhost:8080/register-subscriber
```

- Publish:

```bash
curl -X POST -d '{"publisher_id":"pub1", "topic":"topic1", "payload":"Hello, World!"}' -H "Content-Type: application/json" http://localhost:8080/publish
```
