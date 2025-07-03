package connstore

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
)

func TestStore_Add_Remove_Get(t *testing.T) {
	store := New()

	// Создаем тестовый WebSocket соединение
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()
	}))
	defer server.Close()

	// Подключаемся к тестовому серверу
	wsURL := "ws" + server.URL[4:] // Заменяем http на ws
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}
	defer conn.Close()

	userID := "test-user-123"

	// Тест добавления соединения
	store.Add(userID, conn)

	// Проверяем, что соединение добавлено
	retrievedConn, exists := store.Get(userID)
	if !exists {
		t.Error("Connection should exist after adding")
	}
	if retrievedConn != conn {
		t.Error("Retrieved connection should be the same as added")
	}

	// Тест удаления соединения
	store.Remove(userID)

	// Проверяем, что соединение удалено
	_, exists = store.Get(userID)
	if exists {
		t.Error("Connection should not exist after removing")
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := New()

	// Создаем тестовый WebSocket соединение
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Fatalf("Failed to upgrade connection: %v", err)
		}
		defer conn.Close()
	}))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to test server: %v", err)
	}
	defer conn.Close()

	// Тестируем конкурентный доступ
	done := make(chan bool)

	// Горутина 1: добавляет соединения
	go func() {
		for i := 0; i < 100; i++ {
			userID := fmt.Sprintf("user-%d", i)
			store.Add(userID, conn)
		}
		done <- true
	}()

	// Горутина 2: читает соединения
	go func() {
		for i := 0; i < 100; i++ {
			userID := fmt.Sprintf("user-%d", i)
			store.Get(userID)
		}
		done <- true
	}()

	// Горутина 3: удаляет соединения
	go func() {
		for i := 0; i < 100; i++ {
			userID := fmt.Sprintf("user-%d", i)
			store.Remove(userID)
		}
		done <- true
	}()

	// Ждем завершения всех горутин
	<-done
	<-done
	<-done
}

func TestStore_New(t *testing.T) {
	store := New()

	if store == nil {
		t.Error("New() should return a non-nil store")
	}

	if store.connections == nil {
		t.Error("Store connections map should be initialized")
	}
}
