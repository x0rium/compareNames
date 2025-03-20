package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/x0rium/compareNames/matcher"
)

// RequestBody представляет тело запроса для сравнения имен
type RequestBody struct {
	Name1      string            `json:"name1"`
	Name2      string            `json:"name2"`
	Attributes map[string]bool   `json:"attributes,omitempty"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// SetupRoutes настраивает маршруты API
func SetupRoutes() http.Handler {
	r := mux.NewRouter()
	
	// Middleware для логирования запросов
	r.Use(loggingMiddleware)
	
	// API endpoints
	r.HandleFunc("/api/match_names", handleMatchNames).Methods("POST")
	r.HandleFunc("/health", handleHealth).Methods("GET")
	
	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	
	return c.Handler(r)
}

// loggingMiddleware логирует информацию о запросах
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// handleHealth обрабатывает запросы на проверку работоспособности
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

// handleMatchNames обрабатывает запросы на сравнение имен
func handleMatchNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req RequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}
	
	// Проверка обязательных полей
	if req.Name1 == "" || req.Name2 == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Both name1 and name2 are required"})
		return
	}
	
	// Преобразование атрибутов
	attrs := matcher.CreateAttributes()
	for key, value := range req.Attributes {
		matcher.AddAttribute(attrs, key, value)
	}
	
	// Сравнение имен
	result := matcher.MatchNames(req.Name1, req.Name2, attrs, nil)
	
	// Отправка результата
	json.NewEncoder(w).Encode(result)
}
