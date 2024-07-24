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

func (b *Broker) Subscribe(topic string, subscriber Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
}

func (b *Broker) Publish(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, subscriber := range b.subscribers[msg.Topic] {
		subscriber.Notify(msg)
	}
}

func (b *Broker) RegisterPublisher(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.publishers[id] = true
}

func (b *Broker) IsPublisherRegistered(id string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.publishers[id]
}
