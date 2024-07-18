package main

import (
	"github.com/muneeb-jan/go-broker/internal/messagebroker"
)

func main() {
	broker := messagebroker.NewBroker()

	subscriber1 := &messagebroker.ConcreteSubscriber{ID: "1"}
	subscriber2 := &messagebroker.ConcreteSubscriber{ID: "2"}

	broker.Subscribe("topic1", subscriber1)
	broker.Subscribe("topic1", subscriber2)

	publisher := messagebroker.NewPublisher(broker)
	publisher.Publish("topic1", "Hello, World!")
	publisher.Publish("topic1", "Another message")
}
