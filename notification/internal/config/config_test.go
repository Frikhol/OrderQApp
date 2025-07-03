package config

import (
	"os"
	"testing"
)

func TestGetConfigFromEnv_DefaultValues(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
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

	// Устанавливаем обязательные переменные
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("VERSION", "1.0.0")
	os.Setenv("RABBITMQ_TODO", "test-queue")

	cfg, err := GetConfigFromEnv()
	if err != nil {
		t.Fatalf("GetConfigFromEnv failed: %v", err)
	}

	// Проверяем значения по умолчанию
	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected default log level debug, got %s", cfg.LogLevel)
	}

	if cfg.RabbitMQ.URL != "amqp://guest:guest@rabbitmq:5672/" {
		t.Errorf("Expected default RabbitMQ URL, got %s", cfg.RabbitMQ.URL)
	}

	// Проверяем установленные значения
	if cfg.ServiceName != "test-service" {
		t.Errorf("Expected service name test-service, got %s", cfg.ServiceName)
	}

	if cfg.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", cfg.Version)
	}

	if cfg.RabbitMQ.Query != "test-queue" {
		t.Errorf("Expected RabbitMQ query test-queue, got %s", cfg.RabbitMQ.Query)
	}
}

func TestGetConfigFromEnv_CustomValues(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
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

	// Устанавливаем все переменные
	os.Setenv("SERVICE_NAME", "custom-service")
	os.Setenv("VERSION", "2.0.0")
	os.Setenv("PORT", "9090")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("RABBITMQ_URL", "amqp://user:pass@localhost:5672/")
	os.Setenv("RABBITMQ_TODO", "custom-queue")

	cfg, err := GetConfigFromEnv()
	if err != nil {
		t.Fatalf("GetConfigFromEnv failed: %v", err)
	}

	// Проверяем кастомные значения
	if cfg.ServiceName != "custom-service" {
		t.Errorf("Expected service name custom-service, got %s", cfg.ServiceName)
	}

	if cfg.Version != "2.0.0" {
		t.Errorf("Expected version 2.0.0, got %s", cfg.Version)
	}

	if cfg.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Port)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected log level info, got %s", cfg.LogLevel)
	}

	if cfg.RabbitMQ.URL != "amqp://user:pass@localhost:5672/" {
		t.Errorf("Expected custom RabbitMQ URL, got %s", cfg.RabbitMQ.URL)
	}

	if cfg.RabbitMQ.Query != "custom-queue" {
		t.Errorf("Expected RabbitMQ query custom-queue, got %s", cfg.RabbitMQ.Query)
	}
}

func TestGetConfigFromEnv_MissingRequiredFields(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
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

	// Не устанавливаем обязательные переменные
	_, err := GetConfigFromEnv()
	if err == nil {
		t.Error("Expected error when required fields are missing")
	}
}

func TestConfig_Struct(t *testing.T) {
	cfg := &Config{
		ServiceName: "test",
		Version:     "1.0",
		Port:        "8080",
		LogLevel:    "debug",
		RabbitMQ: RabbitMQ{
			URL:   "amqp://localhost:5672/",
			Query: "test-queue",
		},
	}

	if cfg.ServiceName != "test" {
		t.Errorf("Expected service name test, got %s", cfg.ServiceName)
	}

	if cfg.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", cfg.Version)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Port)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level debug, got %s", cfg.LogLevel)
	}

	if cfg.RabbitMQ.URL != "amqp://localhost:5672/" {
		t.Errorf("Expected RabbitMQ URL, got %s", cfg.RabbitMQ.URL)
	}

	if cfg.RabbitMQ.Query != "test-queue" {
		t.Errorf("Expected RabbitMQ query test-queue, got %s", cfg.RabbitMQ.Query)
	}
}
