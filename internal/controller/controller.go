package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/muneeb-jan/go-broker/internal/messagebroker"
)

type Controller struct {
	broker  *messagebroker.Broker
	devMode bool
	mu      sync.Mutex
}

func NewController(broker *messagebroker.Broker, devMode bool) *Controller {
	return &Controller{
		broker:  broker,
		devMode: devMode,
	}
}

type RegisterPublisherRequest struct {
	ID string `json:"id"`
}

type RegisterPublisherResponse struct {
	Token string `json:"token"`
}

type RegisterSubscriberRequest struct {
	ID       string `json:"id"`
	Topic    string `json:"topic"`
	Listener string `json:"listener"`
}

type RegisterSubscriberResponse struct {
	Token string `json:"token"`
}

type PublishRequest struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}

func (c *Controller) RegisterPublisher(w http.ResponseWriter, r *http.Request) {
	var req RegisterPublisherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.broker.RegisterPublisher(req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := GenerateJWT(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := RegisterPublisherResponse{Token: token}
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) RegisterSubscriber(w http.ResponseWriter, r *http.Request) {
	var req RegisterSubscriberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subscriber := messagebroker.NewConcreteSubscriber(req.ID, req.Listener, c.devMode)
	if err := c.broker.Subscribe(req.Topic, subscriber, req.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := GenerateJWT(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := RegisterSubscriberResponse{Token: token}
	json.NewEncoder(w).Encode(response)
}

func (c *Controller) Publish(w http.ResponseWriter, r *http.Request) {
	var req PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	publisherID := r.Header.Get("User-ID")
	if !c.broker.IsPublisherRegistered(publisherID) {
		http.Error(w, "Publisher not registered", http.StatusForbidden)
		return
	}

	publisher := messagebroker.NewPublisher(c.broker, publisherID)
	publisher.Publish(req.Topic, req.Payload)
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/register-publisher", c.RegisterPublisher)
	mux.HandleFunc("/register-subscriber", c.RegisterSubscriber)
	mux.Handle("/publish", JWTAuthMiddleware(http.HandlerFunc(c.Publish)))
	return mux
}
