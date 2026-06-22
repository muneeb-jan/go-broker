package messagebroker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/muneeb-jan/go-broker/internal/database"
	"github.com/muneeb-jan/go-broker/internal/models"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect test database: %v", err)
	}
	err = db.AutoMigrate(&models.Publisher{}, &models.Subscriber{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	database.DB = db
}

func TestPublishReachesSubscriber(t *testing.T) {
	setupTestDB(t)
	broker := NewBroker()

	var receivedMsg Message
	var mu sync.Mutex

	// Create a test subscriber that captures the message
	testSubscriber := &testSubscriber{
		notifyFunc: func(msg Message) {
			mu.Lock()
			defer mu.Unlock()
			receivedMsg = msg
		},
	}

	// Subscribe to a topic
	subID := "test-sub"
	err := broker.Subscribe("test-topic", testSubscriber, subID)
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish a message
	msg := Message{
		Topic:   "test-topic",
		Payload: map[string]string{"key": "value"},
	}
	count := broker.Publish(msg)

	if count != 1 {
		t.Errorf("Expected 1 subscriber notified, got %d", count)
	}

	// Wait for the async notification
	// Give some time for the goroutine to complete
	select {
	case <-time.After(100 * time.Millisecond):
		// OK, we waited long enough
	}

	mu.Lock()
	defer mu.Unlock()

	if receivedMsg.Topic != "test-topic" {
		t.Errorf("Expected topic 'test-topic', got '%s'", receivedMsg.Topic)
	}

	var payloadMap map[string]string
	payloadBytes, _ := json.Marshal(receivedMsg.Payload)
	json.Unmarshal(payloadBytes, &payloadMap)

	if payloadMap["key"] != "value" {
		t.Errorf("Expected payload key='value', got '%s'", payloadMap["key"])
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	setupTestDB(t)
	broker := NewBroker()

	msg := Message{
		Topic:   "empty-topic",
		Payload: "hello",
	}

	count := broker.Publish(msg)

	if count != 0 {
		t.Errorf("Expected 0 subscribers notified, got %d", count)
	}
}

func TestMultipleSubscribersReceiveMessage(t *testing.T) {
	setupTestDB(t)
	broker := NewBroker()

	var wg sync.WaitGroup
	var mu sync.Mutex
	receivedCount := 0

	// Create 3 subscribers
	for i := 0; i < 3; i++ {
		subID := fmt.Sprintf("sub-%d", i)
		sub := &testSubscriber{
			notifyFunc: func(msg Message) {
				defer wg.Done()
				mu.Lock()
				receivedCount++
				mu.Unlock()
			},
		}
		err := broker.Subscribe("multi-topic", sub, subID)
		if err != nil {
			t.Fatalf("Failed to subscribe %s: %v", subID, err)
		}
		wg.Add(1)
	}

	// Publish
	msg := Message{
		Topic:   "multi-topic",
		Payload: "broadcast",
	}
	count := broker.Publish(msg)

	if count != 3 {
		t.Errorf("Expected 3 subscribers notified, got %d", count)
	}

	// Wait for all notifications
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()

	if receivedCount != 3 {
		t.Errorf("Expected 3 messages received, got %d", receivedCount)
	}
}

func TestPublisherNotFound(t *testing.T) {
	setupTestDB(t)
	broker := NewBroker()

	// Register a publisher
	err := broker.RegisterPublisher("pub1")
	if err != nil {
		t.Fatalf("Failed to register publisher: %v", err)
	}

	// Register a subscriber
	var receivedCount int
	sub := &testSubscriber{
		notifyFunc: func(msg Message) {
			receivedCount++
		},
	}
	err = broker.Subscribe("test-topic", sub, "sub1")
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish as registered publisher
	publisher := NewPublisher(broker, "pub1")
	count := publisher.Publish("test-topic", "hello")
	if count != 1 {
		t.Errorf("Expected 1 subscriber notified, got %d", count)
	}

	// Publish as unregistered publisher
	unregisteredPublisher := NewPublisher(broker, "unknown-pub")
	count = unregisteredPublisher.Publish("test-topic", "hello")
	if count != 0 {
		t.Errorf("Expected 0 subscribers notified for unregistered publisher, got %d", count)
	}
}

func TestSubscriberWithHTTPListener(t *testing.T) {
	setupTestDB(t)
	broker := NewBroker()

	// Create a test HTTP server that receives messages
	var receivedPayload string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var msg struct {
			Topic   string      `json:"topic"`
			Payload interface{} `json:"payload"`
		}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		receivedPayload = fmt.Sprintf("%v", msg.Payload)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Register a publisher
	err := broker.RegisterPublisher("pub1")
	if err != nil {
		t.Fatalf("Failed to register publisher: %v", err)
	}

	// Register a subscriber with HTTP listener
	sub := NewConcreteSubscriber("sub1", server.URL, false)
	err = broker.Subscribe("http-topic", sub, "sub1")
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}

	// Publish a message
	publisher := NewPublisher(broker, "pub1")
	publisher.Publish("http-topic", "hello via http")

	// Wait for HTTP callback
	select {
	case <-time.After(200 * time.Millisecond):
		// OK
	}

	if receivedPayload != "hello via http" {
		t.Errorf("Expected payload 'hello via http', got '%s'", receivedPayload)
	}
}

func TestPublishReturnsSubscriberCount(t *testing.T) {
	setupTestDB(t)
	broker := NewBroker()

	// Register subscribers
	for i := 0; i < 3; i++ {
		sub := &testSubscriber{
			notifyFunc: func(msg Message) {},
		}
		err := broker.Subscribe("topic", sub, fmt.Sprintf("sub-%d", i))
		if err != nil {
			t.Fatalf("Failed to subscribe: %v", err)
		}
	}

	// Publish
	count := broker.Publish(Message{Topic: "topic", Payload: "data"})

	if count != 3 {
		t.Errorf("Expected 3 subscribers notified, got %d", count)
	}

	// Publish to topic with no subscribers
	count = broker.Publish(Message{Topic: "empty-topic", Payload: "data"})
	if count != 0 {
		t.Errorf("Expected 0 subscribers for empty topic, got %d", count)
	}
}

// testSubscriber implements the Subscriber interface for testing
type testSubscriber struct {
	notifyFunc func(Message)
}

func (s *testSubscriber) Notify(msg Message) {
	if s.notifyFunc != nil {
		s.notifyFunc(msg)
	}
}
