package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// PreprocessName предобработка имени: приведение к нижнему регистру, удаление лишних символов
func PreprocessName(name string) string {
	if name == "" {
		return ""
	}

	// Приведение к нижнему регистру
	name = strings.ToLower(name)

	// Преобразование дефисов в пробелы для корректной обработки двойных фамилий
	name = strings.ReplaceAll(name, "-", " ")

	// Обработка апострофов - удаляем их для упрощения сравнения
	name = strings.ReplaceAll(name, "'", "")

	// Удаление лишних пробелов
	re := regexp.MustCompile(`\s+`)
	name = re.ReplaceAllString(name, " ")

	return strings.TrimSpace(name)
}

// NormalizeNameParts разбивает имя на части и возвращает их в нормализованном виде
func NormalizeNameParts(name string) []string {
	// Проверяем, содержит ли имя инициалы
	hasInitials := strings.Contains(name, ".")

	// Если имя содержит инициалы, обрабатываем специальным образом
	if hasInitials {
		// Разбиваем имя на части
		parts := strings.Fields(name)
		result := make([]string, 0, len(parts))

		// Обрабатываем каждую часть
		for _, part := range parts {
			if part != "" {
				// Если это инициал с точкой, сохраняем только букву (без точки)
				if strings.HasSuffix(part, ".") && len(part) == 2 {
					// Добавляем только первую букву инициала
					result = append(result, string(part[0]))
				} else {
					// Иначе добавляем часть как обычно
					result = append(result, part)
				}
			}
		}

		// Если в результате только одна часть, это может быть фамилия с инициалами
		// Например, "Петров И. С." -> ["петров", "и", "с"]
		if len(result) == 1 && len(parts) > 1 {
			// Разбиваем имя на части, заменяя точки на пробелы
			name = strings.ReplaceAll(name, ".", " ")
			parts = strings.Fields(name)
			result = make([]string, 0, len(parts))

			// Обрабатываем каждую часть
			for _, part := range parts {
				if part != "" {
					result = append(result, part)
				}
			}
		}

		return result
	}

	// Если имя не содержит инициалов, обрабатываем как обычно
	parts := strings.Fields(name)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}

	return result
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

// RemoveDiacritics удаляет диакритические знаки из строки
func RemoveDiacritics(s string) string {
	// Преобразуем строку в NFD форму (разделяем символы и диакритические знаки)
	t := []rune(s)
	result := make([]rune, 0, len(t))

	for _, r := range t {
		// Пропускаем диакритические знаки
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		result = append(result, r)
	}

	return string(result)
}

// ContainsRune проверяет, содержит ли строка указанную руну
func ContainsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

// IsLetter проверяет, является ли руна буквой (кириллической или латинской)
func IsLetter(r rune) bool {
	return unicode.IsLetter(r)
}

// IsCyrillicLetter проверяет, является ли руна кириллической буквой
func IsCyrillicLetter(r rune) bool {
	return unicode.Is(unicode.Cyrillic, r)
}

// IsLatinLetter проверяет, является ли руна латинской буквой
func IsLatinLetter(r rune) bool {
	return unicode.Is(unicode.Latin, r)
}

// CountLetters подсчитывает количество букв в строке
func CountLetters(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			count++
		}
	}
	return count
}

// GetFirstLetters возвращает первые буквы каждого слова в строке
func GetFirstLetters(s string) string {
	parts := strings.Fields(s)
	result := make([]rune, 0, len(parts))

	for _, part := range parts {
		if len(part) > 0 {
			result = append(result, []rune(part)[0])
		}
	}

	return string(result)
}

// GetInitials возвращает инициалы из имени
// Например, "Иванов Петр Сергеевич" -> "ИПС"
func GetInitials(name string) string {
	parts := strings.Fields(name)
	result := make([]rune, 0, len(parts))

	for _, part := range parts {
		if len(part) > 0 {
			result = append(result, []rune(part)[0])
		}
	}

	return string(result)
}

// SplitName разбивает полное имя на фамилию, имя и отчество
// Возвращает слайс из трех строк: [фамилия, имя, отчество]
// Если какой-то части нет, соответствующий элемент будет пустой строкой
func SplitName(fullName string) []string {
	parts := strings.Fields(fullName)
	result := make([]string, 3)

	switch len(parts) {
	case 1:
		// Только фамилия или только имя
		result[0] = parts[0]
	case 2:
		// Фамилия и имя
		result[0] = parts[0]
		result[1] = parts[1]
	case 3:
		// Фамилия, имя и отчество
		result[0] = parts[0]
		result[1] = parts[1]
		result[2] = parts[2]
	default:
		// Более трех частей - берем первые три
		if len(parts) > 0 {
			result[0] = parts[0]
		}
		if len(parts) > 1 {
			result[1] = parts[1]
		}
		if len(parts) > 2 {
			result[2] = parts[2]
		}
	}

	return result
}
