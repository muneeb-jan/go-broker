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
	devMode  bool
}

func (cs *ConcreteSubscriber) Notify(msg Message) {
	if cs.devMode {
		fmt.Printf("Subscriber %s received message: %v\n", cs.ID, msg)
	} else {
		payload, _ := json.Marshal(msg)
		_, err := http.Post(cs.Listener, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Printf("Failed to notify subscriber %s: %v\n", cs.ID, err)
		}
	}
}

func NewConcreteSubscriber(id string, listener string, devMode bool) *ConcreteSubscriber {
	return &ConcreteSubscriber{
		ID:       id,
		Listener: listener,
		devMode:  devMode,
	}
}
