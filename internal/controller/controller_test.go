package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/muneeb-jan/go-broker/internal/database"
	"github.com/muneeb-jan/go-broker/internal/messagebroker"
	"github.com/muneeb-jan/go-broker/internal/models"
)

func setupTestController(t *testing.T) (*Controller, *httptest.Server) {
	t.Helper()

	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect test database: %v", err)
	}
	db.AutoMigrate(&models.Publisher{}, &models.Subscriber{})
	database.DB = db

	// Setup broker and controller
	broker := messagebroker.NewBroker()
	ctrl := NewController(broker, true)

	// Setup test server
	server := httptest.NewServer(ctrl.Routes())
	return ctrl, server
}

func TestRegisterPublisher(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	body := `{"id":"pub1"}`
	req, err := http.NewRequest("POST", server.URL+"/register-publisher", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["token"] == "" {
		t.Error("Expected token in response")
	}
}

func TestRegisterSubscriber(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	body := `{"id":"sub1","topic":"test-topic","listener":"http://localhost:8081/listener"}`
	req, err := http.NewRequest("POST", server.URL+"/register-subscriber", bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["token"] == "" {
		t.Error("Expected token in response")
	}
}

func TestPublishWithValidJWT(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	// Register publisher
	pubBody := `{"id":"pub1"}`
	pubReq, _ := http.NewRequest("POST", server.URL+"/register-publisher", bytes.NewBufferString(pubBody))
	pubReq.Header.Set("Content-Type", "application/json")
	pubResp, _ := http.DefaultClient.Do(pubReq)
	defer pubResp.Body.Close()

	var pubResult map[string]string
	json.NewDecoder(pubResp.Body).Decode(&pubResult)
	jwtToken := pubResult["token"]

	// Register subscriber
	subBody := `{"id":"sub1","topic":"test-topic","listener":"http://localhost:8081/listener"}`
	subReq, _ := http.NewRequest("POST", server.URL+"/register-subscriber", bytes.NewBufferString(subBody))
	subReq.Header.Set("Content-Type", "application/json")
	subResp, _ := http.DefaultClient.Do(subReq)
	defer subResp.Body.Close()

	// Publish message
	publishBody := `{"topic":"test-topic","payload":"hello world"}`
	publishReq, _ := http.NewRequest("POST", server.URL+"/publish", bytes.NewBufferString(publishBody))
	publishReq.Header.Set("Content-Type", "application/json")
	publishReq.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(publishReq)
	if err != nil {
		t.Fatalf("Failed to send publish request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}
}

func TestPublishWithBearerPrefix(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	// Register publisher
	pubBody := `{"id":"pub2"}`
	pubReq, _ := http.NewRequest("POST", server.URL+"/register-publisher", bytes.NewBufferString(pubBody))
	pubReq.Header.Set("Content-Type", "application/json")
	pubResp, _ := http.DefaultClient.Do(pubReq)
	defer pubResp.Body.Close()

	var pubResult map[string]string
	json.NewDecoder(pubResp.Body).Decode(&pubResult)
	jwtToken := pubResult["token"]

	// Register subscriber
	subBody := `{"id":"sub2","topic":"test-topic","listener":"http://localhost:8081/listener"}`
	subReq, _ := http.NewRequest("POST", server.URL+"/register-subscriber", bytes.NewBufferString(subBody))
	subReq.Header.Set("Content-Type", "application/json")
	subResp, _ := http.DefaultClient.Do(subReq)
	defer subResp.Body.Close()

	// Publish with Bearer prefix (correct format)
	publishBody := `{"topic":"test-topic","payload":"test"}`
	publishReq, _ := http.NewRequest("POST", server.URL+"/publish", bytes.NewBufferString(publishBody))
	publishReq.Header.Set("Content-Type", "application/json")
	publishReq.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(publishReq)
	if err != nil {
		t.Fatalf("Failed to send publish request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204 with Bearer prefix, got %d", resp.StatusCode)
	}
}

func TestPublishWithoutAuth(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	// Try to publish without auth header
	publishBody := `{"topic":"test-topic","payload":"test"}`
	publishReq, _ := http.NewRequest("POST", server.URL+"/publish", bytes.NewBufferString(publishBody))
	publishReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(publishReq)
	if err != nil {
		t.Fatalf("Failed to send publish request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 without auth, got %d", resp.StatusCode)
	}
}

func TestPublishWithInvalidToken(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	// Try to publish with invalid token
	publishBody := `{"topic":"test-topic","payload":"test"}`
	publishReq, _ := http.NewRequest("POST", server.URL+"/publish", bytes.NewBufferString(publishBody))
	publishReq.Header.Set("Content-Type", "application/json")
	publishReq.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, err := http.DefaultClient.Do(publishReq)
	if err != nil {
		t.Fatalf("Failed to send publish request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 with invalid token, got %d", resp.StatusCode)
	}
}

func TestPublishToTopicWithNoSubscribers(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	// Register publisher
	pubBody := `{"id":"pub3"}`
	pubReq, _ := http.NewRequest("POST", server.URL+"/register-publisher", bytes.NewBufferString(pubBody))
	pubReq.Header.Set("Content-Type", "application/json")
	pubResp, _ := http.DefaultClient.Do(pubReq)
	defer pubResp.Body.Close()

	var pubResult map[string]string
	json.NewDecoder(pubResp.Body).Decode(&pubResult)
	jwtToken := pubResult["token"]

	// Publish to topic with no subscribers
	publishBody := `{"topic":"empty-topic","payload":"test"}`
	publishReq, _ := http.NewRequest("POST", server.URL+"/publish", bytes.NewBufferString(publishBody))
	publishReq.Header.Set("Content-Type", "application/json")
	publishReq.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(publishReq)
	if err != nil {
		t.Fatalf("Failed to send publish request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status 404 for topic with no subscribers, got %d", resp.StatusCode)
	}
}

func TestControllerRoutes(t *testing.T) {
	_, server := setupTestController(t)
	defer server.Close()

	// Test that all routes exist
	routes := []string{"/register-publisher", "/register-subscriber", "/publish"}
	for _, route := range routes {
		req, _ := http.NewRequest("GET", server.URL+route, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Failed to access route %s: %v", route, err)
		}
		resp.Body.Close()

		if route == "/publish" {
			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Route %s: expected 401, got %d", route, resp.StatusCode)
			}
		} else {
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("Route %s: expected 400, got %d", route, resp.StatusCode)
			}
		}
	}
}
