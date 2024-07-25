# go-broker

Simple Message broker written in Go.

## Installation

Runnning the message broker is as easy as any other go app. Just use the command

```bash
go run cmd/go-broker/main.go
```

Want to develop this code further? Say no more. It can be rather hectic to have a subscriber with a listener API always running. The solution? Running the app in development mode. How? Just add the dev flag in the run command.

```bash
go run cmd/go-broker/main.go --dev
```

In this mode, the app prints all the published messages instead of posting it to listener. For example:

```bash
Subscriber sub1 received message: {topic1 Hello, World!}
```

## APIs

### POST APIs

- Register Publisher: 

```bash
curl -X POST -d '{"id":"pub1"}' -H "Content-Type: application/json" http://localhost:8080/register-publisher
```

- Register Subscriber:

```bash
curl -X POST -d '{"id":"sub1", "topic":"topic1", "listener":"http://localhost:8081/listener"}' -H "Content-Type: application/json" http://localhost:8080/register-subscriber
```

- Publish:

```bash
curl -X POST -d '{"publisher_id":"pub1", "topic":"topic1", "payload":"Hello, World!"}' -H "Content-Type: application/json" http://localhost:8080/publish
```
