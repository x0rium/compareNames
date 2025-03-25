package compare

import (
	"strings"
	"time"

	"github.com/x0rium/compareNames/matcher/translit"
)

// MatchResult содержит результаты сравнения имен
type MatchResult struct {
	ExactMatch                bool    // Точное совпадение
	Score                     int     // Оценка совпадения (0-100)
	MatchType                 string  // Тип совпадения ("match", "possible_match", "no_match")
	BestMatch1                string  // Лучшее совпадение для первого имени
	BestMatch2                string  // Лучшее совпадение для второго имени
	LevenshteinScore          float64 // Оценка по Левенштейну
	JaroWinklerScore          float64 // Оценка по Джаро-Винклеру
	PhoneticScore             float64 // Фонетическая оценка
	DoubleMetaphoneScore      float64 // Оценка по Double Metaphone
	CosineScore               float64 // Оценка по косинусному сходству
	AdditionalAttributesScore float64 // Оценка по дополнительным атрибутам
	ProcessingTimeMS          int64   // Время обработки в миллисекундах
	FromCache                 bool    // Результат получен из кэша
}

// MatchAttributes карта дополнительных атрибутов
type MatchAttributes map[string]struct{ Match bool }

// NameMatcher интерфейс для матчера имён
type NameMatcher interface {
	// Получение вариаций имён
	GetNameVariations(name string) []string

	// Предобработка имени
	PreprocessName(name string) string

	// Проверка на наличие инициалов
	HasInitials(name string) bool

	// Нормализация частей имени
	NormalizeNameParts(name string) []string

	// Получение ключа кэша (если используется кэширование)
	GetCacheKey(name1, name2 string, attrs MatchAttributes) string

	// Методы для получения конфигурационных параметров
	GetMinExactMatchScore() int
	GetMinPossibleMatchScore() int
	GetConfig() ConfigProvider
}

// MatchNames основной метод для сравнения имен
func MatchNames(name1, name2 string, attrs MatchAttributes, matcher NameMatcher) MatchResult {
	startTime := time.Now()

	// 1. Проверяем на инициалы
	hasInitials1 := strings.Contains(name1, ".") || matcher.HasInitials(name1)
	hasInitials2 := strings.Contains(name2, ".") || matcher.HasInitials(name2)

	// Если одно из имен содержит инициалы, пробуем обработать этот случай
	if hasInitials1 || hasInitials2 {
		initialsResult, found := ProcessInitials(name1, name2, hasInitials1, hasInitials2)
		if found {
			// Заполняем недостающие поля в результате
			return MatchResult{
				ExactMatch:           initialsResult.IsMatch && initialsResult.Score >= 95,
				Score:                initialsResult.Score,
				MatchType:            initialsResult.MatchType,
				BestMatch1:           initialsResult.BestMatch1,
				BestMatch2:           initialsResult.BestMatch2,
				LevenshteinScore:     0.9, // Примерные значения для инициалов
				JaroWinklerScore:     0.9,
				PhoneticScore:        initialsResult.PhoneticScore,
				DoubleMetaphoneScore: 1.0,
				ProcessingTimeMS:     time.Since(startTime).Milliseconds(),
			}
		}
	}

	// 2. Предобработка имен
	processedName1 := matcher.PreprocessName(name1)
	processedName2 := matcher.PreprocessName(name2)

	// Проверка на пустые имена
	if processedName1 == "" || processedName2 == "" {
		return MatchResult{
			ExactMatch:           false,
			Score:                0,
			MatchType:            "no_match",
			LevenshteinScore:     0.0,
			JaroWinklerScore:     0.0,
			PhoneticScore:        0.0,
			DoubleMetaphoneScore: 0.0,
			ProcessingTimeMS:     time.Since(startTime).Milliseconds(),
		}
	}

	// 3. Проверка на точное совпадение после предобработки
	if processedName1 == processedName2 {
		return MatchResult{
			ExactMatch:           true,
			Score:                100,
			MatchType:            "match",
			LevenshteinScore:     1.0,
			JaroWinklerScore:     1.0,
			PhoneticScore:        1.0,
			DoubleMetaphoneScore: 1.0,
			ProcessingTimeMS:     time.Since(startTime).Milliseconds(),
		}
	}

	// 4. Проверка на разные языки (кириллица/латиница)
	isName1Cyrillic := translit.IsCyrillic(processedName1)
	isName2Cyrillic := translit.IsCyrillic(processedName2)

	// 5. Выбираем стратегию сравнения в зависимости от языка имен
	if isName1Cyrillic != isName2Cyrillic {
		// Сравнение имен на разных алфавитах
		differentAlphabetsResult := CompareDifferentAlphabets(
			processedName1,
			processedName2,
			isName1Cyrillic,
			startTime,
			matcher.GetNameVariations,
			matcher.NormalizeNameParts,
			func(parts1, parts2 []string) float64 {
				return CompareNameParts(parts1, parts2, matcher.GetConfig())
			},
			matcher.GetMinExactMatchScore(),
		)

		// Если не нашли хорошего совпадения, используем сравнение как для одного алфавита
		if differentAlphabetsResult.MatchType == "no_match" {
			return convertSameAlphabetResult(
				CompareSameAlphabet(processedName1, processedName2, attrs, matcher.GetConfig(), startTime),
			)
		}

		// Иначе возвращаем результат сравнения разных алфавитов
		return convertDifferentAlphabetsResult(differentAlphabetsResult)
	}

	// 6. Если оба имени на одном языке, используем стандартный алгоритм
	sameAlphabetResult := CompareSameAlphabet(processedName1, processedName2, attrs, matcher.GetConfig(), startTime)
	return convertSameAlphabetResult(sameAlphabetResult)
}

// convertSameAlphabetResult преобразует SameAlphabetResult в MatchResult
func convertSameAlphabetResult(result SameAlphabetResult) MatchResult {
	return MatchResult{
		ExactMatch:                result.ExactMatch,
		Score:                     result.Score,
		MatchType:                 result.MatchType,
		BestMatch1:                result.BestMatch1,
		BestMatch2:                result.BestMatch2,
		LevenshteinScore:          result.LevenshteinScore,
		JaroWinklerScore:          result.JaroWinklerScore,
		PhoneticScore:             result.PhoneticScore,
		DoubleMetaphoneScore:      result.DoubleMetaphoneScore,
		CosineScore:               result.CosineScore,
		AdditionalAttributesScore: result.AdditionalAttributesScore,
		ProcessingTimeMS:          result.ProcessingTime,
	}
}

// convertDifferentAlphabetsResult преобразует DifferentAlphabetsResult в MatchResult
func convertDifferentAlphabetsResult(result DifferentAlphabetsResult) MatchResult {
	return MatchResult{
		ExactMatch:           result.ExactMatch,
		Score:                result.Score,
		MatchType:            result.MatchType,
		BestMatch1:           result.BestMatch1,
		BestMatch2:           result.BestMatch2,
		LevenshteinScore:     result.LevenshteinScore,
		JaroWinklerScore:     result.JaroWinklerScore,
		PhoneticScore:        result.PhoneticScore,
		DoubleMetaphoneScore: result.DoubleMetaphoneScore,
		ProcessingTimeMS:     result.ProcessingTimeMS,
	}
}

// Пустой комментарий, чтобы убрать заглушки
