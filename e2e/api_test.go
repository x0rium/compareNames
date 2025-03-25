package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/x0rium/compareNames/matcher"
)

// TestCase представляет собой тестовый случай для API
type TestCase struct {
	Name               string `json:"name"`
	Name1              string `json:"name1"`
	Name2              string `json:"name2"`
	ExpectedScore      int    `json:"expectedScore"`
	ExpectedMatchType  string `json:"expectedMatchType"`
	ExpectedExactMatch bool   `json:"expectedExactMatch"`
}

// RequestBody представляет структуру запроса к API
type RequestBody struct {
	Name1 string `json:"name1"`
	Name2 string `json:"name2"`
}

// LoadTestCases загружает тестовые случаи из файла JSON
func LoadTestCases(t *testing.T) []TestCase {
	// Находим путь к текущему исполняемому файлу
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Не удалось определить путь к текущему файлу")
	}

	// Получаем директорию, в которой находится тестовый файл
	dir := filepath.Dir(filename)

	// Формируем полный путь к файлу cases.json
	casesPath := filepath.Join(dir, "cases.json")

	// Открываем файл с тестовыми случаями
	file, err := os.Open(casesPath)
	if err != nil {
		t.Fatalf("Ошибка при открытии файла с тестовыми случаями: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("Ошибка при чтении файла с тестовыми случаями: %v", err)
	}

	// Декодируем JSON в структуру
	var testCases []TestCase
	err = json.Unmarshal(content, &testCases)
	if err != nil {
		t.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	return testCases
}

// TestMatchNames тестирует API сравнения имен
func TestMatchNames(t *testing.T) {
	// Устанавливаем тестовый сервер
	setupTestServer(t)
	defer teardownTestServer(t)

	// Загружаем тестовые случаи
	testCases := LoadTestCases(t)

	// Создаем HTTP-клиент
	client := &http.Client{}

	// URL для API сравнения имен
	apiURL := fmt.Sprintf("%s/api/match_names", baseURL)

	// Выполняем тесты для каждого случая
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Создаем тело запроса
			requestBody := RequestBody{
				Name1: tc.Name1,
				Name2: tc.Name2,
			}

			// Логируем информацию о тестовом случае
			t.Logf("📝 Тест: %s", tc.Name)
			t.Logf("📥 Запрос: name1=\"%s\", name2=\"%s\"", tc.Name1, tc.Name2)
			t.Logf("🎯 Ожидаем: score=%d, matchType=%s, exactMatch=%v",
				tc.ExpectedScore, tc.ExpectedMatchType, tc.ExpectedExactMatch)

			// Сериализуем тело запроса в JSON
			requestJSON, err := json.Marshal(requestBody)
			if err != nil {
				t.Fatalf("Ошибка при сериализации запроса: %v", err)
			}

			// Создаем HTTP-запрос
			req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestJSON))
			if err != nil {
				t.Fatalf("Ошибка при создании запроса: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Отправляем запрос
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Ошибка при отправке запроса: %v", err)
			}
			defer resp.Body.Close()

			// Проверяем успешный код ответа
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("Неожиданный код ответа: %d", resp.StatusCode)
			}

			// Читаем ответ
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Ошибка при чтении ответа: %v", err)
			}

			// Декодируем ответ
			var matchResult matcher.MatchResult
			err = json.Unmarshal(respBody, &matchResult)
			if err != nil {
				t.Fatalf("Ошибка при декодировании ответа: %v", err)
			}

			// Логируем полученный результат
			t.Logf("📤 Результат: score=%d, matchType=%s, exactMatch=%v",
				matchResult.Score, matchResult.MatchType, matchResult.ExactMatch)
			t.Logf("📊 Метрики: Левенштейн=%.2f, Джаро-Винклер=%.2f, Фонетика=%.2f, DoubleMetaphone=%.2f",
				matchResult.LevenshteinScore, matchResult.JaroWinklerScore,
				matchResult.PhoneticScore, matchResult.DoubleMetaphoneScore)

			// Проверяем попадание score в нужный диапазон в соответствии с matchType
			scoreRangeMatches := false
			if tc.ExpectedMatchType == "exact_match" && matchResult.Score == 100 {
				// Для exact_match ожидаем ровно 100
				scoreRangeMatches = true
			} else if tc.ExpectedMatchType == "match" && matchResult.Score > 90 {
				// Для match ожидаем > 90
				scoreRangeMatches = true
			} else if tc.ExpectedMatchType == "possible_match" && matchResult.Score >= 70 && matchResult.Score <= 90 {
				// Для possible_match ожидаем 70-90
				scoreRangeMatches = true
			} else if tc.ExpectedMatchType == "no_match" && matchResult.Score < 70 {
				// Для no_match ожидаем < 70
				scoreRangeMatches = true
			}

			// Определяем статус проверки
			passStatus := "✅ PASS"
			if !scoreRangeMatches ||
				matchResult.MatchType != tc.ExpectedMatchType ||
				matchResult.ExactMatch != tc.ExpectedExactMatch {
				passStatus = "❌ FAIL"
			}
			t.Logf("%s: %s <-> %s", passStatus, tc.Name1, tc.Name2)

			// Проверяем результаты
			if !scoreRangeMatches {
				t.Logf("⚠️ Оценка %d не соответствует диапазону для %s",
					matchResult.Score, tc.ExpectedMatchType)

				// Выводим ожидаемый диапазон в зависимости от типа совпадения
				if tc.ExpectedMatchType == "exact_match" {
					t.Logf("   Ожидаемый диапазон: = 100")
				} else if tc.ExpectedMatchType == "match" {
					t.Logf("   Ожидаемый диапазон: > 90")
				} else if tc.ExpectedMatchType == "possible_match" {
					t.Logf("   Ожидаемый диапазон: 70-90")
				} else if tc.ExpectedMatchType == "no_match" {
					t.Logf("   Ожидаемый диапазон: < 70")
				}
			}

			if matchResult.MatchType != tc.ExpectedMatchType {
				t.Errorf("❌ Ожидаемый тип совпадения: %s, получен: %s", tc.ExpectedMatchType, matchResult.MatchType)
			}

			if matchResult.ExactMatch != tc.ExpectedExactMatch {
				t.Errorf("❌ Ожидаемое точное совпадение: %v, получено: %v", tc.ExpectedExactMatch, matchResult.ExactMatch)
			}
		})
	}
}

// TestHealthCheck проверяет endpoint проверки работоспособности
func TestHealthCheck(t *testing.T) {
	// Устанавливаем тестовый сервер
	setupTestServer(t)
	defer teardownTestServer(t)

	// URL для проверки работоспособности
	healthURL := fmt.Sprintf("%s/health", baseURL)

	// Отправляем запрос
	resp, err := http.Get(healthURL)
	if err != nil {
		t.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем успешный код ответа
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Неожиданный код ответа: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Ошибка при чтении ответа: %v", err)
	}

	// Проверяем, что содержимое ответа содержит информацию о статусе
	if !strings.Contains(string(body), `"status"`) {
		t.Errorf("Неожиданный формат ответа: %s", body)
	}
}
