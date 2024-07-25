package messagebroker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ConcreteSubscriber struct {
	ID       string
	Listener string
}

func (cs *ConcreteSubscriber) Notify(msg Message) {
	fmt.Printf("Subscriber %s received message: %v\n", cs.ID, msg)
	payload, _ := json.Marshal(msg)
	http.Post(cs.Listener, "application/json", bytes.NewBuffer(payload))
}
