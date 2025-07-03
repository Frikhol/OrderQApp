package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"notification_service/internal/config"
	"notification_service/internal/infra/auth"
	"notification_service/proto/auth_service"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MockAuthService реализует mock gRPC сервис аутентификации
type MockAuthService struct {
	auth_service.UnimplementedAuthServiceServer
	validTokens map[string]string // token -> userID
}

func (m *MockAuthService) ValidateToken(ctx context.Context, req *auth_service.ValidateTokenRequest) (*auth_service.ValidateTokenResponse, error) {
	userID, exists := m.validTokens[req.Token]
	if !exists {
		return &auth_service.ValidateTokenResponse{Success: false}, nil
	}
	return &auth_service.ValidateTokenResponse{Success: true, UserId: userID}, nil
}

func TestNotificationService_Integration(t *testing.T) {
	// Настройка тестового окружения
	originalEnv := make(map[string]string)
	for _, key := range []string{"SERVICE_NAME", "VERSION", "PORT", "LOG_LEVEL", "RABBITMQ_URL", "RABBITMQ_TODO"} {
		if val := os.Getenv(key); val != "" {
			originalEnv[key] = val
		}
	}

	// Очищаем переменные окружения для теста
	for _, key := range []string{"SERVICE_NAME", "VERSION", "PORT", "LOG_LEVEL", "RABBITMQ_URL", "RABBITMQ_TODO"} {
		os.Unsetenv(key)
	}

	// Восстанавливаем переменные окружения после теста
	defer func() {
		for key, val := range originalEnv {
			os.Setenv(key, val)
		}
	}()

	// Устанавливаем тестовые переменные окружения
	os.Setenv("SERVICE_NAME", "test-notification-service")
	os.Setenv("VERSION", "1.0.0")
	os.Setenv("PORT", "0") // Используем порт 0 для автоматического выбора свободного порта
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	os.Setenv("RABBITMQ_TODO", "test-notifications")

	// Получаем конфигурацию
	cfg, err := config.GetConfigFromEnv()
	if err != nil {
		t.Fatalf("Failed to get config: %v", err)
	}

	// Создаем mock gRPC сервер аутентификации
	mockAuthServer := grpc.NewServer()
	mockAuthService := &MockAuthService{
		validTokens: map[string]string{
			"valid-token-1": "user-1",
			"valid-token-2": "user-2",
		},
	}
	auth_service.RegisterAuthServiceServer(mockAuthServer, mockAuthService)

	// Запускаем mock gRPC сервер
	grpcListener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create gRPC listener: %v", err)
	}
	defer grpcListener.Close()

	go func() {
		if err := mockAuthServer.Serve(grpcListener); err != nil {
			t.Errorf("Failed to serve gRPC: %v", err)
		}
	}()

	// Подключаемся к mock gRPC серверу
	grpcAddr := grpcListener.Addr().String()
	authConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer authConn.Close()

	// Создаем auth клиент
	authClient := auth.NewAuthClient(auth_service.NewAuthServiceClient(authConn))

	// Тестируем валидацию токенов
	t.Run("TokenValidation", func(t *testing.T) {
		// Тест валидного токена
		userID, err := authClient.ValidateToken("valid-token-1")
		if err != nil {
			t.Errorf("Expected valid token to pass validation: %v", err)
		}
		if userID != "user-1" {
			t.Errorf("Expected user ID 'user-1', got '%s'", userID)
		}

		// Тест невалидного токена
		_, err = authClient.ValidateToken("invalid-token")
		if err == nil {
			t.Error("Expected invalid token to fail validation")
		}
	})

	// Тестируем WebSocket соединения
	t.Run("WebSocketConnections", func(t *testing.T) {
		// Создаем тестовый HTTP сервер с WebSocket обработчиком
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Здесь должен быть ваш WebSocket обработчик
			// Для простоты тестируем только HTTP ответ
			if r.URL.Path == "/ws" {
				if authHeader := r.Header.Get("Authorization"); authHeader != "" {
					if strings.HasPrefix(authHeader, "Bearer ") {
						token := strings.TrimPrefix(authHeader, "Bearer ")
						if _, exists := mockAuthService.validTokens[token]; exists {
							w.WriteHeader(http.StatusSwitchingProtocols)
							return
						}
					}
				}
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// Тест успешного WebSocket соединения
		t.Run("SuccessfulConnection", func(t *testing.T) {
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
			header := http.Header{}
			header.Set("Authorization", "Bearer valid-token-1")

			conn, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
			if err != nil {
				t.Fatalf("Failed to connect to WebSocket: %v", err)
			}
			defer conn.Close()

			if resp.StatusCode != http.StatusSwitchingProtocols {
				t.Errorf("Expected status 101, got %d", resp.StatusCode)
			}
		})

		// Тест неуспешного WebSocket соединения
		t.Run("FailedConnection", func(t *testing.T) {
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
			header := http.Header{}
			header.Set("Authorization", "Bearer invalid-token")

			_, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
			if err == nil {
				t.Error("Expected connection to fail with invalid token")
			}
			if resp != nil && resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("Expected status 401, got %d", resp.StatusCode)
			}
		})
	})

	// Тестируем конфигурацию
	t.Run("Configuration", func(t *testing.T) {
		if cfg.ServiceName != "test-notification-service" {
			t.Errorf("Expected service name 'test-notification-service', got '%s'", cfg.ServiceName)
		}

		if cfg.Version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got '%s'", cfg.Version)
		}

		if cfg.LogLevel != "debug" {
			t.Errorf("Expected log level 'debug', got '%s'", cfg.LogLevel)
		}

		if cfg.RabbitMQ.Query != "test-notifications" {
			t.Errorf("Expected RabbitMQ query 'test-notifications', got '%s'", cfg.RabbitMQ.Query)
		}
	})
}

// Benchmark тесты для производительности
func BenchmarkAuthClient_ValidateToken(b *testing.B) {
	// Настройка benchmark
	authClient := &auth.AuthClient{} // Здесь нужен mock клиент

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = authClient.ValidateToken("test-token")
	}
}

func BenchmarkWebSocketConnection(b *testing.B) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusSwitchingProtocols)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			b.Fatalf("Failed to connect: %v", err)
		}
		conn.Close()
	}
}
