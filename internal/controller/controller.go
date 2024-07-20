package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/muneeb-jan/go-broker/internal/messagebroker"
)

type Controller struct {
	broker *messagebroker.Broker
	mu     sync.Mutex
}

func NewController(broker *messagebroker.Broker) *Controller {
	return &Controller{
		broker: broker,
	}
}

type PublishRequest struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

type SubscribeRequest struct {
	Topic string `json:"topic"`
	ID    string `json:"id"`
}

func (c *Controller) Publish(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publisher := messagebroker.NewPublisher(c.broker)
	publisher.Publish(req.Topic, req.Payload)

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscriber := &messagebroker.ConcreteSubscriber{ID: req.ID}
	c.broker.Subscribe(req.Topic, subscriber)

	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/publish", c.Publish)
	mux.HandleFunc("/subscribe", c.Subscribe)
	return mux
}
