package matcher

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// MatchNames сравнивает два имени с указанной конфигурацией
// Экспортированная функция для использования в других пакетах
func MatchNames(name1, name2 string, attrs Attributes, cfg *Config) MatchResult {
	
	// Создаем экземпляр NameMatcher с указанной конфигурацией
	matcher := NewNameMatcher(cfg)
	
	// Выполняем сравнение имен
	result := matcher.MatchNames(name1, name2, attrs)
	
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
