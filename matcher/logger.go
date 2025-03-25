package matcher

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// LogEntry структура для записи в лог
type LogEntry struct {
	Timestamp      time.Time  `json:"timestamp"`
	Name1          string     `json:"name1"`
	Name2          string     `json:"name2"`
	Score          int        `json:"score"`
	MatchType      string     `json:"match_type"`
	Metrics        LogMetrics `json:"metrics"`
	ProcessingTime int64      `json:"processing_time_ms"`
	Attributes     Attributes `json:"attributes,omitempty"`
}

// LogMetrics метрики для логирования
type LogMetrics struct {
	Levenshtein     float64 `json:"levenshtein"`
	JaroWinkler     float64 `json:"jaro_winkler"`
	Phonetic        float64 `json:"phonetic"`
	DoubleMetaphone float64 `json:"double_metaphone"`
	Cosine          float64 `json:"cosine"`
}

// LogPossibleMatch логирует сомнительные совпадения
func LogPossibleMatch(name1, name2 string, attrs Attributes, result MatchResult) {
	// Проверяем переменную окружения GO_TEST
	if os.Getenv("GO_TEST") == "1" {
		return
	}

	// Логирование подробной информации для анализа
	log.Printf("Сомнительное совпадение обнаружено: %s <-> %s (Оценка: %d, Тип: %s)",
		name1, name2, result.Score, result.MatchType)
	log.Printf("Метрики: Левенштейн=%.2f, Джаро-Винклер=%.2f, Фонетическая=%.2f, DoubleMetaphone=%.2f",
		result.LevenshteinScore, result.JaroWinklerScore, result.PhoneticScore, result.DoubleMetaphoneScore)

	entry := LogEntry{
		Timestamp: time.Now(),
		Name1:     name1,
		Name2:     name2,
		Score:     result.Score,
		MatchType: result.MatchType,
		Metrics: LogMetrics{
			Levenshtein:     result.LevenshteinScore,
			JaroWinkler:     result.JaroWinklerScore,
			Phonetic:        result.PhoneticScore,
			DoubleMetaphone: result.DoubleMetaphoneScore,
			Cosine:          result.CosineScore,
		},
		ProcessingTime: result.ProcessingTimeMS,
		Attributes:     attrs,
	}

	// Сохраняем в JSON файл
	logJSON, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		log.Printf("Ошибка маршалинга лога: %v", err)
		return
	}

	// Открываем файл лога
	f, err := os.OpenFile("name_matching.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Ошибка открытия файла лога: %v", err)
		return
	}
	defer f.Close()

	// Записываем в файл
	if _, err := f.Write(append(logJSON, '\n')); err != nil {
		log.Printf("Ошибка записи в файл лога: %v", err)
	}
}
