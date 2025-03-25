package e2e

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/x0rium/compareNames/api"
)

var (
	testServer *http.Server
	serverPort int
	baseURL    string
)

// setupTestServer запускает тестовый сервер на случайном порту
func setupTestServer(t *testing.T) {
	// Выбираем случайный порт в диапазоне 10000-60000
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	serverPort = 10000 + r.Intn(50000)
	baseURL = fmt.Sprintf("http://localhost:%d", serverPort)

	// Настраиваем роуты
	router := api.SetupRoutes()

	// Создаем сервер
	testServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Printf("Запуск тестового сервера на порту %d", serverPort)
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Ошибка запуска тестового сервера: %v", err)
		}
	}()

	// Даем серверу время для запуска
	time.Sleep(100 * time.Millisecond)
}

// teardownTestServer останавливает тестовый сервер
func teardownTestServer(t *testing.T) {
	if testServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := testServer.Shutdown(ctx); err != nil {
			t.Errorf("Ошибка при остановке тестового сервера: %v", err)
		}

		log.Println("Тестовый сервер остановлен")
	}
}

// TestMain выполняется перед всеми тестами
func TestMain(m *testing.M) {
	// Настраиваем тестовую среду
	log.Println("Настройка тестовой среды...")

	// Запускаем тесты
	code := m.Run()

	// Завершаем тесты
	log.Println("Тесты завершены")

	// Возвращаем статус код
	os.Exit(code)
}
