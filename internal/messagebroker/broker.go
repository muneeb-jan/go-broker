package messagebroker

import (
	"sync"

	"github.com/muneeb-jan/go-broker/internal/database"
	"github.com/muneeb-jan/go-broker/internal/models"
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

// Subscribe registers a subscriber to a topic and saves the subscription to the database
func (b *Broker) Subscribe(topic string, subscriber Subscriber, subscriberID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Save subscriber to the database
	dbSubscriber := &models.Subscriber{
		ID:    subscriberID,
		Topic: topic,
	}
	if err := database.DB.Create(dbSubscriber).Error; err != nil {
		return err
	}

	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
	return nil
}

// Publish sends a message to all subscribers of a topic concurrently
func (b *Broker) Publish(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, subscriber := range b.subscribers[msg.Topic] {
		go subscriber.Notify(msg) // Notify each subscriber in a separate Go routine
	}
}

// RegisterPublisher adds a publisher to the broker and saves it to the database
func (b *Broker) RegisterPublisher(id string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Save publisher to the database
	dbPublisher := &models.Publisher{
		ID: id,
	}
	if err := database.DB.Create(dbPublisher).Error; err != nil {
		return err
	}

	b.publishers[id] = true
	return nil
}

// IsPublisherRegistered checks if a publisher is registered in memory or the database
func (b *Broker) IsPublisherRegistered(id string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.publishers[id] {
		return true
	}

	// Check if the publisher exists in the database
	var publisher models.Publisher
	if err := database.DB.Where("id = ?", id).First(&publisher).Error; err == nil {
		// Cache the publisher in memory
		b.mu.Lock()
		b.publishers[id] = true
		b.mu.Unlock()
		return true
	}

	return false
}
