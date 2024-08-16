package messagebroker

import (
	"sync"
)

type Message struct {
	Topic   string
	Payload interface{}
}

type Subscriber interface {
	Notify(msg Message)
}

type Broker struct {
	subscribers map[string][]Subscriber
	publishers  map[string]bool
	mu          sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subscribers: make(map[string][]Subscriber),
		publishers:  make(map[string]bool),
	}
}

// Subscribe registers a subscriber to a topic
func (b *Broker) Subscribe(topic string, subscriber Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
}

// Publish sends a message to all subscribers of a topic concurrently
func (b *Broker) Publish(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, subscriber := range b.subscribers[msg.Topic] {
		go subscriber.Notify(msg) // Notify each subscriber in a separate Go routine
	}
}

// RegisterPublisher adds a publisher to the broker
func (b *Broker) RegisterPublisher(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.publishers[id] = true
}

// IsPublisherRegistered checks if a publisher is registered
func (b *Broker) IsPublisherRegistered(id string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.publishers[id]
}
