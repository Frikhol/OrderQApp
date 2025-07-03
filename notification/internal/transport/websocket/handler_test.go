package websocket

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"notification_service/internal/connstore"

	"github.com/gorilla/websocket"
)

// MockAuth реализует интерфейс auth.Auth для тестирования
type MockAuth struct {
	shouldFail bool
	userID     string
}

func (m *MockAuth) ValidateToken(token string) (string, error) {
	if m.shouldFail {
		return "", errors.New("invalid token")
	}
	if token == "" {
		return "", errors.New("no token provided")
	}
	return m.userID, nil
}

func TestHandler_HandleUser_NoToken(t *testing.T) {
	store := connstore.New()
	mockAuth := &MockAuth{shouldFail: false, userID: "test-user"}
	handler := NewHandler(mockAuth, store)

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handler.HandleUser))
	defer server.Close()

	// Делаем запрос без токена
	resp, err := http.Get(server.URL + "/ws")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Ожидаем ошибку 400 Bad Request
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_HandleUser_InvalidToken(t *testing.T) {
	store := connstore.New()
	mockAuth := &MockAuth{shouldFail: true, userID: ""}
	handler := NewHandler(mockAuth, store)

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handler.HandleUser))
	defer server.Close()

	// Делаем запрос с неверным токеном
	req, err := http.NewRequest("GET", server.URL+"/ws", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer invalid-token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Ожидаем ошибку 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestHandler_HandleUser_ValidToken(t *testing.T) {
	store := connstore.New()
	mockAuth := &MockAuth{shouldFail: false, userID: "test-user-123"}
	handler := NewHandler(mockAuth, store)

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handler.HandleUser))
	defer server.Close()

	// Подключаемся через WebSocket с валидным токеном
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	// Создаем заголовки с токеном
	header := http.Header{}
	header.Set("Authorization", "Bearer valid-token")

	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer conn.Close()

	// Проверяем, что соединение установлено успешно
	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("Expected status 101, got %d", resp.StatusCode)
	}

	// Проверяем, что пользователь добавлен в store
	_, exists := store.Get("test-user-123")
	if !exists {
		t.Error("User should be added to store after successful connection")
	}
}

func TestHandler_HandleUser_EmptyBearerToken(t *testing.T) {
	store := connstore.New()
	mockAuth := &MockAuth{shouldFail: false, userID: "test-user"}
	handler := NewHandler(mockAuth, store)

	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(handler.HandleUser))
	defer server.Close()

	// Делаем запрос с пустым токеном после "Bearer "
	req, err := http.NewRequest("GET", server.URL+"/ws", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer ")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Ожидаем ошибку 400 Bad Request
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
}

func TestHandler_NewHandler(t *testing.T) {
	store := connstore.New()
	mockAuth := &MockAuth{shouldFail: false, userID: "test-user"}

	handler := NewHandler(mockAuth, store)

	if handler == nil {
		t.Error("NewHandler should return a non-nil handler")
	}

	if handler.auth != mockAuth {
		t.Error("Handler should have the provided auth client")
	}

	if handler.store != store {
		t.Error("Handler should have the provided store")
	}

	if handler.upgrader == nil {
		t.Error("Handler should have an initialized upgrader")
	}
}
