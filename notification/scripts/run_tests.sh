#!/bin/bash

# Скрипт для запуска тестов сервиса уведомлений

set -e

echo "🧪 Запуск тестов сервиса уведомлений..."

# Переходим в директорию сервиса
cd "$(dirname "$0")/.."

# Проверяем, что Go установлен
if ! command -v go &> /dev/null; then
    echo "❌ Go не установлен. Пожалуйста, установите Go."
    exit 1
fi

# Очищаем кэш модулей
echo "📦 Очистка кэша модулей..."
go clean -modcache

# Загружаем зависимости
echo "📥 Загрузка зависимостей..."
go mod download

# Запускаем unit тесты
echo "🔬 Запуск unit тестов..."
go test -v ./internal/connstore/...
go test -v ./internal/config/...
go test -v ./internal/transport/websocket/...

# Запускаем интеграционные тесты
echo "🔗 Запуск интеграционных тестов..."
go test -v -tags=integration ./...

# Запускаем benchmark тесты
echo "⚡ Запуск benchmark тестов..."
go test -bench=. -benchmem ./...

# Проверяем покрытие кода
echo "📊 Проверка покрытия кода..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "✅ Все тесты завершены!"
echo "📄 Отчет о покрытии сохранен в coverage.html"

# Показываем статистику покрытия
echo "📈 Статистика покрытия:"
go tool cover -func=coverage.out | tail -1 