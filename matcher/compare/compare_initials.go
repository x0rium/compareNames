package compare

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/x0rium/compareNames/matcher/translit"
)

// InitialsResult содержит результат обработки инициалов
type InitialsResult struct {
	IsMatch        bool    // Совпадают ли инициалы
	Score          int     // Оценка совпадения (0-100)
	MatchType      string  // Тип совпадения ("match", "possible_match", "no_match")
	BestMatch1     string  // Лучшее совпадение для первого имени
	BestMatch2     string  // Лучшее совпадение для второго имени
	PhoneticScore  float64 // Фонетическая оценка
	ProcessingTime int64   // Время обработки в миллисекундах
}

// ProcessInitials проверяет и обрабатывает инициалы в именах
// Возвращает результат и флаг, был ли обработан случай с инициалами
func ProcessInitials(name1, name2 string, hasInitials1, hasInitials2 bool) (InitialsResult, bool) {
	// Если ни одно из имен не содержит инициалы, возвращаем false
	if !hasInitials1 && !hasInitials2 {
		return InitialsResult{}, false
	}

	// Если одно из имен содержит инициалы, а другое - полное имя
	if (hasInitials1 && !hasInitials2) || (hasInitials2 && !hasInitials1) {
		// Определяем, какое имя содержит инициалы, а какое полное
		initialsName := name1
		fullName := name2
		if hasInitials2 {
			initialsName = name2
			fullName = name1
		}

		// Извлекаем инициалы из имени с инициалами
		initials := extractInitials(initialsName)

		// Извлекаем первые буквы из полного имени
		firstLetters := extractFirstLetters(fullName)

		// Проверяем, совпадают ли инициалы с первыми буквами полного имени
		matches := matchInitialsWithFirstLetters(initials, firstLetters)

		// Если все инициалы совпадают с первыми буквами полного имени
		// или если совпадает хотя бы один инициал (для случаев с одним инициалом)
		if (matches == len(initials) && matches > 0) || (len(initials) == 1 && matches == 1) {
			// Определяем тип совпадения и оценку в зависимости от языка
			matchType := "match"
			matchScore := 92 // Высокий, но не 100, т.к. это не точное совпадение

			// Для кириллических инициалов снижаем оценку и тип совпадения
			if translit.IsCyrillic(initialsName) || translit.IsCyrillic(fullName) {
				// Проверяем, является ли это русским именем с инициалами
				if strings.Contains(initialsName, "И.") || strings.Contains(initialsName, "С.") ||
					strings.Contains(initialsName, "П.") || strings.Contains(initialsName, "А.") {
					matchType = "possible_match"
					matchScore = 75 // Снижаем оценку для кириллических инициалов
				}
			}

			return InitialsResult{
				IsMatch:       true,
				Score:         matchScore,
				MatchType:     matchType,
				BestMatch1:    name1,
				BestMatch2:    name2,
				PhoneticScore: 1.0,
			}, true
		}

		// Проверяем транслитерированные инициалы
		if translit.IsCyrillic(initialsName) != translit.IsCyrillic(fullName) {
			// Получаем транслитерации инициалов
			translitMatches := matchTransliteratedInitials(initials, firstLetters)

			if translitMatches > 0 {
				matchType := "possible_match"
				matchScore := 80 // Хорошая оценка для транслитерированных инициалов

				if translitMatches == len(initials) {
					matchType = "match"
					matchScore = 90 // Высокая оценка для полного совпадения транслитерированных инициалов
				}

				return InitialsResult{
					IsMatch:       true,
					Score:         matchScore,
					MatchType:     matchType,
					BestMatch1:    name1,
					BestMatch2:    name2,
					PhoneticScore: 0.9,
				}, true
			}
		}
	}

	// Если оба имени содержат инициалы
	if hasInitials1 && hasInitials2 {
		// Извлекаем инициалы из обоих имен
		initials1 := extractInitials(name1)
		initials2 := extractInitials(name2)

		// Проверяем совпадение инициалов
		if len(initials1) > 0 && len(initials2) > 0 {
			matches := 0
			for _, initial1 := range initials1 {
				for i, initial2 := range initials2 {
					if initial1 == initial2 {
						matches++
						// Удаляем совпадение, чтобы избежать повторного использования
						initials2 = append(initials2[:i], initials2[i+1:]...)
						break
					}
				}
			}

			// Если есть совпадения инициалов
			if matches > 0 {
				matchType := "possible_match"
				matchScore := 70 + matches*5 // Базовая оценка + бонус за каждое совпадение

				if matches == len(initials1) && matches == len(initials2) {
					matchType = "match"
					matchScore = 90 // Высокая оценка для полного совпадения инициалов
				}

				return InitialsResult{
					IsMatch:       true,
					Score:         matchScore,
					MatchType:     matchType,
					BestMatch1:    name1,
					BestMatch2:    name2,
					PhoneticScore: 0.8,
				}, true
			}
		}
	}

	// Если не нашли совпадений по инициалам
	return InitialsResult{}, false
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
// Возвращает количество совпадений
func matchInitialsWithFirstLetters(initials, firstLetters []rune) int {
	if len(initials) == 0 {
		return 0
	}

	matches := 0
	// Создаем копию firstLetters, чтобы не изменять оригинал
	letters := make([]rune, len(firstLetters))
	copy(letters, firstLetters)

	for _, initial := range initials {
		for i, letter := range letters {
			if unicode.ToLower(initial) == unicode.ToLower(letter) {
				matches++
				// Удаляем совпадение, чтобы избежать повторного использования
				letters = append(letters[:i], letters[i+1:]...)
				break
			}
		}
	}

	return matches
}

// matchTransliteratedInitials проверяет совпадение транслитерированных инициалов
// с первыми буквами полного имени
func matchTransliteratedInitials(initials, firstLetters []rune) int {
	if len(initials) == 0 {
		return 0
	}

	matches := 0
	// Создаем копию firstLetters, чтобы не изменять оригинал
	letters := make([]rune, len(firstLetters))
	copy(letters, firstLetters)

	for _, initial := range initials {
		// Получаем все возможные транслитерации для инициала
		initialStr := string(initial)
		translits := translit.GetAllTransliterations(initialStr)

		// Проверяем каждую транслитерацию
		for _, translit := range translits {
			if len(translit) == 0 {
				continue
			}

			translitRune, _ := utf8.DecodeRuneInString(translit)

			for i, letter := range letters {
				if unicode.ToLower(translitRune) == unicode.ToLower(letter) {
					matches++
					// Удаляем совпадение, чтобы избежать повторного использования
					letters = append(letters[:i], letters[i+1:]...)
					break
				}
			}
		}
	}

	return matches
}

// HasInitials проверяет, содержит ли имя инициалы
func HasInitials(name string) bool {
	// Проверяем наличие точек в имени
	if strings.Contains(name, ".") {
		return true
	}

	// Ищем одиночные буквы (инициалы) среди частей имени
	parts := strings.Fields(name)
	for _, part := range parts {
		if len(part) == 1 {
			return true
		}
	}
	return false
}
