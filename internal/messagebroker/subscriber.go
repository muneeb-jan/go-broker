package messagebroker

import "fmt"

type ConcreteSubscriber struct {
	ID string
}

func (cs *ConcreteSubscriber) Notify(msg Message) {
	fmt.Printf("Subscriber %s received message: %v\n", cs.ID, msg)
}
