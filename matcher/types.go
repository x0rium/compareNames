package matcher

import (
	"sync"
	"time"
)

// Attribute представляет дополнительный атрибут для сравнения
type Attribute struct {
	Match bool `json:"match"`
}

// Attributes карта дополнительных атрибутов
type Attributes map[string]Attribute

// NameMatcher основной тип для сравнения имен
type NameMatcher struct {
	Config          Config
	nameVariantions map[string][]string // Кэш для хранения вариаций имен
	cache           *Cache              // Кэш для хранения результатов сравнения
	mutex           sync.RWMutex        // Мьютекс для потокобезопасности
}

// MatchResult содержит результаты сравнения имен
type MatchResult struct {
	ExactMatch                bool    `json:"exact_match"`
	Score                     int     `json:"score"`
	MatchType                 string  `json:"match_type"`
	BestMatch1                string  `json:"best_match1,omitempty"`
	BestMatch2                string  `json:"best_match2,omitempty"`
	LevenshteinScore          float64 `json:"levenshtein_score,omitempty"`
	JaroWinklerScore          float64 `json:"jaro_winkler_score,omitempty"`
	PhoneticScore             float64 `json:"phonetic_score,omitempty"`
	DoubleMetaphoneScore      float64 `json:"double_metaphone_score,omitempty"`
	CosineScore               float64 `json:"cosine_score,omitempty"`
	AdditionalAttributesScore float64 `json:"additional_attributes_score,omitempty"`
	ProcessingTimeMS          int64   `json:"processing_time_ms"`
	FromCache                 bool    `json:"from_cache,omitempty"`
}

// NameMatchMetrics структура с метриками совпадения
type NameMatchMetrics struct {
	LevenshteinScore     float64
	JaroWinklerScore     float64
	PhoneticScore        float64
	DoubleMetaphoneScore float64
	CosineScore          float64
	AttributesScore      float64
}

// Cache представляет кэш для результатов сравнения
type Cache struct {
	items   map[string]CacheItem
	keys    []string
	maxSize int
	mutex   sync.RWMutex
	TTL     time.Duration // Время жизни элемента кэша
}

// CacheItem представляет элемент кэша с временем создания
type CacheItem struct {
	Result     MatchResult
	CreateTime time.Time
}

// NewNameMatcher создает новый экземпляр NameMatcher
func NewNameMatcher(cfg *Config) *NameMatcher {
	matcher := &NameMatcher{
		nameVariantions: make(map[string][]string),
	}

	// Используем конфигурацию по умолчанию, если не предоставлена
	if cfg == nil {
		matcher.Config = DefaultConfig()
	} else {
		matcher.Config = *cfg

		// Если в конфигурации не указан размер n-граммы, устанавливаем по умолчанию
		if matcher.Config.NGramSize <= 0 {
			matcher.Config.NGramSize = 3
		}

		// Если не указан максимальный размер кэша, устанавливаем по умолчанию
		if matcher.Config.MaxCacheSize <= 0 {
			matcher.Config.MaxCacheSize = CacheSize
		}
	}

	// Инициализируем кэш, если включено кэширование
	if matcher.Config.EnableCaching {
		matcher.cache = NewCache(matcher.Config.MaxCacheSize, 15*time.Minute) // TTL 15 минут
	}

	return matcher
}

// Вспомогательная функция для вычисления абсолютного значения разницы
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
