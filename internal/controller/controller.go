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

type RegisterPublisherRequest struct {
	ID string `json:"id"`
}

type RegisterSubscriberRequest struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Listener string `json:"listener"`
}

type PublishRequest struct {
	PublisherID string      `json:"publisher_id"`
	Topic       string      `json:"topic"`
	Payload     interface{} `json:"payload"`
}

func (c *Controller) RegisterPublisher(w http.ResponseWriter, r *http.Request) {
	var req RegisterPublisherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.broker.RegisterPublisher(req.ID)
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) RegisterSubscriber(w http.ResponseWriter, r *http.Request) {
	var req RegisterSubscriberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscriber := &messagebroker.ConcreteSubscriber{ID: req.ID, Listener: req.Listener}
	c.broker.Subscribe(req.Topic, subscriber)
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Publish(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publisher := messagebroker.NewPublisher(c.broker, req.PublisherID)
	publisher.Publish(req.Topic, req.Payload)
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register-publisher", c.RegisterPublisher)
	mux.HandleFunc("/register-subscriber", c.RegisterSubscriber)
	mux.HandleFunc("/publish", c.Publish)
	return mux
}
