package matcher

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"

	"github.com/x0rium/compareNames/matcher/similarity"
	"github.com/x0rium/compareNames/matcher/translit"
)

// MatchNames сравнивает два имени с указанной конфигурацией
// Экспортированная функция для использования в других пакетах
func MatchNames(name1, name2 string, attrs Attributes, cfg *Config) MatchResult {
	// Создаем экземпляр NameMatcher с указанной конфигурацией
	_ = NewNameMatcher(cfg)

	// Инициализируем результат
	var result MatchResult

	// Проверяем точное совпадение
	if strings.EqualFold(name1, name2) {
		result.ExactMatch = true
		result.Score = 100
		result.MatchType = "exact_match"
		return result
	}

	// Если имена не совпадают точно, выполняем расширенное сравнение
	bestLevenshteinScore := 0.0
	bestJaroWinklerScore := 0.0

	// Разбиваем имена на части для учета возможных перестановок
	name1Parts := strings.Fields(name1)
	name2Parts := strings.Fields(name2)

	// Генерируем все возможные перестановки частей для имен
	name1Permutations := []string{name1}
	name2Permutations := []string{name2}

	// Если имя состоит из нескольких частей, генерируем перестановки
	if len(name1Parts) > 1 {
		for i := 0; i < len(name1Parts); i++ {
			for j := i + 1; j < len(name1Parts); j++ {
				// Переставляем части местами
				permParts := make([]string, len(name1Parts))
				copy(permParts, name1Parts)
				permParts[i], permParts[j] = permParts[j], permParts[i]
				name1Permutations = append(name1Permutations, strings.Join(permParts, " "))
			}
		}
	}

	if len(name2Parts) > 1 {
		for i := 0; i < len(name2Parts); i++ {
			for j := i + 1; j < len(name2Parts); j++ {
				// Переставляем части местами
				permParts := make([]string, len(name2Parts))
				copy(permParts, name2Parts)
				permParts[i], permParts[j] = permParts[j], permParts[i]
				name2Permutations = append(name2Permutations, strings.Join(permParts, " "))
			}
		}
	}

	// Объединяем перестановки с вариантами транслитерации
	var allName1Variants []string
	var allName2Variants []string

	// Получаем все варианты транслитерации для всех перестановок имен
	for _, perm := range name1Permutations {
		variants := translit.GetAllTransliterations(perm)
		allName1Variants = append(allName1Variants, variants...)
	}

	for _, perm := range name2Permutations {
		variants := translit.GetAllTransliterations(perm)
		allName2Variants = append(allName2Variants, variants...)
	}

	// Сравниваем каждую пару вариантов и выбираем наилучший результат
	for _, variant1 := range allName1Variants {
		for _, variant2 := range allName2Variants {
			// Вычисляем расстояние Левенштейна
			currentLevenshteinDist := levenshteinDistance(strings.ToLower(variant1), strings.ToLower(variant2))
			maxLen := math.Max(float64(len(variant1)), float64(len(variant2)))
			currentLevenshteinScore := 0.0
			if maxLen > 0 {
				currentLevenshteinScore = 1.0 - float64(currentLevenshteinDist)/maxLen
			}

			// Вычисляем сходство Джаро-Винклера
			currentJaroWinklerScore := jaroWinklerSimilarity(variant1, variant2)

			// Обновляем наилучшие оценки
			if currentLevenshteinScore > bestLevenshteinScore {
				bestLevenshteinScore = currentLevenshteinScore
			}
			if currentJaroWinklerScore > bestJaroWinklerScore {
				bestJaroWinklerScore = currentJaroWinklerScore
			}
		}
	}

	// Вычисляем фонетические оценки
	phoneticScore := similarity.SoundexSimilarity(name1, name2)
	doubleMetaphoneScore := similarity.DoubleMetaphoneSimilarity(name1, name2)

	// Устанавливаем оценки в результат (округляем до двух знаков после запятой)
	result.LevenshteinScore = math.Round(bestLevenshteinScore*100) / 100
	result.JaroWinklerScore = math.Round(bestJaroWinklerScore*100) / 100
	result.PhoneticScore = math.Round(phoneticScore*100) / 100
	result.DoubleMetaphoneScore = math.Round(doubleMetaphoneScore*100) / 100

	// Устанавливаем конфигурацию, если не предоставлена
	if cfg == nil {
		defaultCfg := DefaultConfig()
		cfg = &defaultCfg
	}

	// Вычисляем базовую оценку как взвешенное среднее всех оценок
	avgScore := (result.LevenshteinScore*cfg.LevenshteinWeight +
		result.JaroWinklerScore*cfg.JaroWinklerWeight +
		result.PhoneticScore*cfg.PhoneticWeight +
		result.DoubleMetaphoneScore*cfg.DoubleMetaphoneWeight)

	// Применяем специальные бонусы в соответствии с бизнес-требованиями
	baseScore := avgScore // Сохраняем базовую оценку
	totalBonus := 0.0     // Суммарный бонус

	// Бонус 1: Транслитерация между алфавитами (до 12%)
	if translit.IsCyrillic(name1) != translit.IsCyrillic(name2) {
		// Для транслитерации даём бонус (больше для хороших фонетических совпадений)
		if result.PhoneticScore > 0.8 {
			totalBonus += 0.12 // 12% бонус для хороших фонетических совпадений
		} else {
			totalBonus += 0.08 // 8% бонус для обычных случаев
		}
	}

	// Бонус 2: Перестановки частей ФИО (до 12%)
	// Проверяем, является ли одно имя перестановкой другого
	if isNamePartsPermutation(name1, name2) {
		totalBonus += 0.12 // 12% бонус
	}

	// Бонус 3: Обработка инициалов (до 15%)
	if hasInitialsAtStart(name1, name2) {
		if strings.Contains(name1, ".") || strings.Contains(name2, ".") {
			totalBonus += 0.10 // 10% бонус за инициалы с точками
		} else {
			totalBonus += 0.15 // 15% бонус за инициалы без точек (полный инициал)
		}
	}

	// Бонус 4: Обработка дефисных имен (до 8%)
	if hasHyphenatedName(name1) || hasHyphenatedName(name2) {
		totalBonus += 0.08 // 8% бонус
	}

	// Бонус 5: Обработка уменьшительных/альтернативных форм имен (до 12%)
	if isNameFormVariation(name1, name2) {
		totalBonus += 0.12 // 12% бонус
	}

	// Применяем совокупный бонус к базовой оценке, но не больше 30%
	if totalBonus > 0.3 {
		totalBonus = 0.3
	}

	// Применяем бонус к базовой оценке
	avgScore = baseScore * (1.0 + totalBonus)

	// Ограничиваем максимальное значение до 0.99 (чтобы оставить 100% только для точных совпадений)
	if avgScore > 0.99 {
		avgScore = 0.99
	}

	// Переводим в шкалу 0-100
	result.Score = int(math.Round(avgScore * 100))

	// Определяем тип совпадения на основе оценки
	if result.Score >= cfg.MatchThreshold {
		result.MatchType = "match"
	} else if result.Score >= cfg.PossibleMatchThreshold {
		result.MatchType = "possible_match"
	} else {
		result.MatchType = "no_match"
	}

	// Логируем сомнительные совпадения для дальнейшего анализа
	if result.MatchType == "possible_match" {
		LogPossibleMatch(name1, name2, attrs, result)
	}

	return result
}

// hasInitialsAtStart проверяет, начинается ли одно из имен с инициалов, а другое с полных имен
func hasInitialsAtStart(name1, name2 string) bool {

	// Проверяем, содержит ли name1 инициалы, а name2 - полные имена
	if containsInitials(name1) && !containsInitials(name2) {
		initials := extractInitials(name1)
		firstLetters := extractFirstLetters(name2)

		// Проверяем, совпадают ли инициалы с первыми буквами полного имени
		return matchInitialsWithFirstLetters(initials, firstLetters)
	}

	// Проверяем, содержит ли name2 инициалы, а name1 - полные имена
	if containsInitials(name2) && !containsInitials(name1) {
		initials := extractInitials(name2)
		firstLetters := extractFirstLetters(name1)

		// Проверяем, совпадают ли инициалы с первыми буквами полного имени
		return matchInitialsWithFirstLetters(initials, firstLetters)
	}

	return false
}

// containsInitials проверяет, содержит ли имя инициалы
func containsInitials(name string) bool {
	// Проверяем наличие точек (характерно для инициалов)
	if strings.Contains(name, ".") {
		return true
	}

	// Проверяем наличие одиночных букв (могут быть инициалами без точек)
	parts := strings.Fields(name)
	for _, part := range parts {
		if len(part) == 1 {
			return true
		}
	}

	return false
}

// extractInitials извлекает инициалы из имени
func extractInitials(name string) []rune {
	var initials []rune
	parts := strings.Fields(name)

	for _, part := range parts {
		if len(part) == 1 {
			// Одиночная буква (инициал без точки)
			r, _ := utf8.DecodeRuneInString(part)
			initials = append(initials, r)
		} else if len(part) == 2 && strings.HasSuffix(part, ".") {
			// Инициал с точкой (например, "И.")
			r, _ := utf8.DecodeRuneInString(part)
			initials = append(initials, r)
		}
	}

	return initials
}

// extractFirstLetters извлекает первые буквы из полного имени
func extractFirstLetters(name string) []rune {
	var firstLetters []rune
	parts := strings.Fields(name)

	for _, part := range parts {
		if len(part) > 0 {
			r, _ := utf8.DecodeRuneInString(part)
			firstLetters = append(firstLetters, r)
		}
	}

	return firstLetters
}

// matchInitialsWithFirstLetters проверяет, совпадают ли инициалы с первыми буквами полного имени
func matchInitialsWithFirstLetters(initials, firstLetters []rune) bool {
	if len(initials) == 0 {
		return false
	}

	matches := 0
	for _, initial := range initials {
		for i, letter := range firstLetters {
			if initial == letter {
				matches++
				// Удаляем совпадение, чтобы избежать повторного использования
				firstLetters = append(firstLetters[:i], firstLetters[i+1:]...)
				break
			}
		}
	}

	// Если все инициалы совпадают с первыми буквами полного имени
	return matches == len(initials)
}

// MatchNamesWithDefaultConfig сравнивает два имени с конфигурацией по умолчанию
func MatchNamesWithDefaultConfig(name1, name2 string) MatchResult {
	return MatchNames(name1, name2, nil, nil)
}

// MatchNamesWithAttributes сравнивает два имени с дополнительными атрибутами
func MatchNamesWithAttributes(name1, name2 string, attrs Attributes) MatchResult {
	return MatchNames(name1, name2, attrs, nil)
}

// CreateAttribute создает атрибут для сравнения
func CreateAttribute(match bool) Attribute {
	return Attribute{Match: match}
}

// CreateAttributes создает карту атрибутов для сравнения
func CreateAttributes() Attributes {
	return make(Attributes)
}

// AddAttribute добавляет атрибут в карту атрибутов
func AddAttribute(attrs Attributes, name string, match bool) {
	attrs[name] = Attribute{Match: match}
}

// isNamePartsPermutation проверяет, является ли одно имя перестановкой другого
func isNamePartsPermutation(name1, name2 string) bool {
	// Разбиваем имена на части
	parts1 := strings.Fields(name1)
	parts2 := strings.Fields(name2)

	// Если разное количество частей, это не перестановка
	if len(parts1) != len(parts2) {
		return false
	}

	// Сортировка не поможет при разных регистрах и транслитерации
	// Поэтому проверяем, что каждая часть name1 имеет соответствие в name2
	// Используем более сложный алгоритм сравнения для учета транслитерации

	matched := make([]bool, len(parts2))

	for _, part1 := range parts1 {
		found := false

		for i, part2 := range parts2 {
			if !matched[i] {
				// Проверяем не только точное совпадение, но и фонетическое
				if strings.EqualFold(part1, part2) ||
					similarity.SoundexMatch(part1, part2) ||
					similarity.RussianSoundexMatch(part1, part2) {
					matched[i] = true
					found = true
					break
				}
			}
		}

		if !found {
			return false
		}
	}

	// Все части нашли соответствие
	return true
}

// hasHyphenatedName проверяет, содержит ли имя дефис
func hasHyphenatedName(name string) bool {
	return strings.Contains(name, "-")
}

// isNameFormVariation проверяет, является ли одно имя вариацией другого (например, "Александр" и "Саша")
func isNameFormVariation(name1, name2 string) bool {
	// Словарь соответствий имен и их уменьшительных/альтернативных форм
	nameVariations := map[string][]string{
		"александр": {"саша", "шура", "алекс", "саня"},
		"алексей":   {"леша", "леха", "алеша", "alex", "alexey", "aleksei"},
		"анастасия": {"настя", "ася", "стася"},
		"анна":      {"аня", "анюта", "анечка"},
		"дмитрий":   {"дима", "димуля", "митя", "dmitry", "dmitri"},
		"екатерина": {"катя", "катерина", "катюша"},
		"елена":     {"лена", "леночка", "helen", "helena"},
		"иван":      {"ваня", "иванушка", "ivan"},
		"мария":     {"маша", "машенька", "mary", "maria"},
		"михаил":    {"миша", "michael", "mikhail"},
		"ольга":     {"оля", "оленька", "olga", "olya"},
		"сергей":    {"серега", "сергеич", "sergey", "sergei"},
		"татьяна":   {"таня", "танюша", "tatiana", "tanya"},
		"юрий":      {"юра", "yuri", "yury", "jurij", "juri"},
		// Можно добавить больше соответствий
	}

	// Преобразуем имена к нижнему регистру для сравнения
	lowerName1 := strings.ToLower(name1)
	lowerName2 := strings.ToLower(name2)

	// Разбиваем имена на слова
	words1 := strings.Fields(lowerName1)
	words2 := strings.Fields(lowerName2)

	// Проверяем, что хотя бы одно слово является вариацией другого
	for _, word1 := range words1 {
		for _, word2 := range words2 {
			// Сначала проверяем, равны ли слова напрямую
			if word1 == word2 {
				return true
			}

			// Затем проверяем, является ли одно слово вариацией другого
			if variations, ok := nameVariations[word1]; ok {
				for _, variation := range variations {
					if variation == word2 {
						return true
					}
				}
			}

			if variations, ok := nameVariations[word2]; ok {
				for _, variation := range variations {
					if variation == word1 {
						return true
					}
				}
			}
		}
	}

	return false
}

// PrintMatchResult выводит результат сравнения в консоль
func PrintMatchResult(result MatchResult) {
	fmt.Printf("Результат сравнения:\n")
	fmt.Printf("  Точное совпадение: %v\n", result.ExactMatch)
	fmt.Printf("  Оценка: %d\n", result.Score)
	fmt.Printf("  Тип совпадения: %s\n", result.MatchType)

	if result.BestMatch1 != "" && result.BestMatch2 != "" {
		fmt.Printf("  Лучшее совпадение 1: %s\n", result.BestMatch1)
		fmt.Printf("  Лучшее совпадение 2: %s\n", result.BestMatch2)
	}

	fmt.Printf("  Оценки алгоритмов:\n")
	fmt.Printf("    Левенштейн: %.4f\n", result.LevenshteinScore)
	fmt.Printf("    Джаро-Винклер: %.4f\n", result.JaroWinklerScore)
	fmt.Printf("    Фонетическая: %.4f\n", result.PhoneticScore)
	fmt.Printf("    Double Metaphone: %.4f\n", result.DoubleMetaphoneScore)
	fmt.Printf("    Косинусная: %.4f\n", result.CosineScore)

	if result.AdditionalAttributesScore > 0 {
		fmt.Printf("    Дополнительные атрибуты: %.4f\n", result.AdditionalAttributesScore)
	}

	fmt.Printf("  Время обработки: %d мс\n", result.ProcessingTimeMS)

	if result.FromCache {
		fmt.Printf("  Результат получен из кэша\n")
	}
}
