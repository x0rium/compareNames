package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/x0rium/compareNames/matcher"
	"github.com/x0rium/compareNames/middleware"
)

// RequestBody структура для запроса к API
type RequestBody struct {
	Name1        string             `json:"name1"`
	Name2        string             `json:"name2"`
	Attributes   matcher.Attributes `json:"attributes,omitempty"`
	Config       *matcher.Config    `json:"config,omitempty"`
	DisableCache bool               `json:"disable_cache,omitempty"`
}

// ErrorResponse структура для ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// MatchNamesHandler обработчик для /api/match_names
func MatchNamesHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	// Парсим тело запроса
	var requestBody RequestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		sendErrorResponse(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if requestBody.Name1 == "" || requestBody.Name2 == "" {
		sendErrorResponse(w, "Both name1 and name2 are required", http.StatusBadRequest)
		return
	}

	// Настраиваем конфигурацию
	config := requestBody.Config
	if config == nil {
		defaultConfig := matcher.DefaultConfig()
		config = &defaultConfig
	}

	// Устанавливаем кэширование, если указано в запросе
	if requestBody.DisableCache {
		config.EnableCaching = false
	}

	// Выполняем сравнение имен
	result := matcher.MatchNames(
		requestBody.Name1,
		requestBody.Name2,
		requestBody.Attributes,
		config,
	)

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding response: %v", err)
		sendErrorResponse(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// HealthCheckHandler обработчик для проверки работоспособности API
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

// Вспомогательная функция для отправки ответа с ошибкой
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Printf("Error encoding error response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// SetupRoutes настраивает маршруты для API
func SetupRoutes() http.Handler {
	router := mux.NewRouter()

	// API endpoint для сравнения имен
	router.HandleFunc("/api/match_names", MatchNamesHandler).Methods("POST")

	// Endpoint для проверки работоспособности API
	router.HandleFunc("/health", HealthCheckHandler).Methods("GET")

	// Применяем middleware для логирования
	return middleware.Logging(router)
}
