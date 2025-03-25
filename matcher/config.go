package matcher

// Константы для настройки алгоритма
const (
	MinExactMatchScore    = 90   // Минимальный балл для точного совпадения
	MinPossibleMatchScore = 70   // Минимальный балл для возможного совпадения
	MaxProcessTimeMS      = 100  // Максимальное время обработки в миллисекундах
	CacheSize             = 1000 // Максимальный размер кэша результатов
)

// Config структура с настройками для алгоритма сравнения имен
type Config struct {
	// Веса для алгоритмов сравнения
	LevenshteinWeight     float64 `json:"levenshtein_weight"`
	JaroWinklerWeight     float64 `json:"jaro_winkler_weight"`
	PhoneticWeight        float64 `json:"phonetic_weight"`
	DoubleMetaphoneWeight float64 `json:"double_metaphone_weight"`
	CosineWeight          float64 `json:"cosine_weight"`
	AdditionalAttrsWeight float64 `json:"additional_attributes_weight"`

	// Пороговые значения
	JaroWinklerThreshold   float64 `json:"jaro_winkler_threshold"`
	LevenshteinPrefixScale float64 `json:"levenshtein_prefix_scale"`
	ExactMatchThreshold    int     `json:"exact_match_threshold"`
	MatchThreshold         int     `json:"match_threshold"`
	PossibleMatchThreshold int     `json:"possible_match_threshold"`

	// Параметры транслитерации
	EnableTransliteration    bool     `json:"enable_transliteration"`
	TransliterationStandards []string `json:"transliteration_standards"`

	// Параметры перестановки
	EnableNamePartPermutation bool `json:"enable_name_part_permutation"`

	// Другие параметры
	NGramSize     int  `json:"ngram_size"`
	EnableCaching bool `json:"enable_caching"`
	MaxCacheSize  int  `json:"max_cache_size"`
	EnableLogging bool `json:"enable_logging"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		// Веса для алгоритмов - сумма = 1.0
		LevenshteinWeight:     0.2,
		JaroWinklerWeight:     0.3,
		PhoneticWeight:        0.3,
		DoubleMetaphoneWeight: 0.2,
		CosineWeight:          0.0, // Не используется в базовой версии
		AdditionalAttrsWeight: 0.0, // Не используется в базовой версии

		// Пороговые значения
		JaroWinklerThreshold:   0.85,
		LevenshteinPrefixScale: 0.1,
		ExactMatchThreshold:    100, // 100% для точного совпадения
		MatchThreshold:         90,  // 90% для совпадения (в соответствии с БТ)
		PossibleMatchThreshold: 70,  // 70% для возможного совпадения (в соответствии с БТ)

		// Параметры транслитерации
		EnableTransliteration:    true,
		TransliterationStandards: []string{"gost", "iso9", "bgnpcgn", "ungegn"},

		// Параметры перестановки
		EnableNamePartPermutation: true,

		// Другие параметры
		NGramSize:     3,
		EnableCaching: true,
		MaxCacheSize:  CacheSize,
		EnableLogging: true,
	}
}
