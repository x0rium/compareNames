package similarity

import (
	"strings"
	"unicode"
)

// DoubleMetaphoneSimilarity вычисляет схожесть двух строк по алгоритму Double Metaphone
func DoubleMetaphoneSimilarity(s1, s2 string) float64 {
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Разбиваем строки на слова
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	// Если нет слов - возвращаем 0
	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Считаем совпадения Double Metaphone кодов
	matchCount := 0
	totalWords := max(len(words1), len(words2))

	// Для каждого слова в первой строке
	for _, word1 := range words1 {
		primary1, secondary1 := DoubleMetaphone(word1)

		// Ищем совпадение во второй строке
		for _, word2 := range words2 {
			primary2, secondary2 := DoubleMetaphone(word2)

			// Проверяем совпадение по первичным и вторичным кодам
			if (primary1 == primary2 || primary1 == secondary2 ||
				secondary1 == primary2 || secondary1 == secondary2) &&
				primary1 != "" && primary2 != "" {
				matchCount++
				break
			}
		}
	}

	// Возвращаем нормализованную оценку
	return float64(matchCount) / float64(totalWords)
}

// DoubleMetaphone реализует алгоритм Double Metaphone
// Возвращает два кода - первичный и вторичный
func DoubleMetaphone(word string) (string, string) {
	// Приводим к верхнему регистру и удаляем не-буквенные символы
	word = cleanString(word)
	if word == "" {
		return "", ""
	}

	// Инициализируем переменные
	length := len(word)
	current := 0
	primary := ""
	secondary := ""

	// Обрабатываем начальные звуки
	if length > 1 {
		// Особые случаи для начала слова
		if word[0:2] == "KN" || word[0:2] == "GN" || word[0:2] == "PN" || word[0:2] == "AE" || word[0:2] == "WR" {
			current = 1
		}

		// Начальная 'X' превращается в 'S'
		if word[0] == 'X' {
			primary += "S"
			secondary += "S"
			current = 1
		}

		// Начальная 'WH' превращается в 'W'
		if word[0:2] == "WH" {
			primary += "W"
			secondary += "W"
			current = 2
		}
	}

	// Основной цикл обработки
	for current < length && (len(primary) < 4 || len(secondary) < 4) {
		c := word[current]

		// Пропускаем двойные согласные, кроме 'C'
		if c == 'C' && current > 0 && word[current-1] == 'C' {
			current++
			continue
		}

		// Обработка в зависимости от символа
		switch c {
		case 'A', 'E', 'I', 'O', 'U', 'Y':
			if current == 0 {
				// Начальные гласные сохраняются
				primary += "A"
				secondary += "A"
			}
			current++

		case 'B':
			// 'B' -> 'P' если в конце слова после 'M'
			if current > 0 && word[current-1] == 'M' && current+1 >= length {
				primary += "P"
				secondary += "P"
			} else {
				primary += "P"
				secondary += "P"
			}

			// Пропускаем 'B' если она в конце слова после 'M'
			if current+1 < length && word[current+1] == 'B' {
				current += 2
			} else {
				current++
			}

		case 'C':
			// Различные правила для 'C'
			if current > 0 && word[current-1] == 'S' && current+1 < length &&
				(word[current+1] == 'I' || word[current+1] == 'E' || word[current+1] == 'Y') {
				// 'SCI', 'SCE', 'SCY' -> пропускаем 'C'
				current++
			} else if current+2 < length && word[current+1] == 'I' && word[current+2] == 'A' {
				// 'CIA' -> 'X'
				primary += "X"
				secondary += "X"
				current += 3
			} else if current+1 < length &&
				(word[current+1] == 'I' || word[current+1] == 'E' || word[current+1] == 'Y') {
				// 'CI', 'CE', 'CY' -> 'S'
				primary += "S"
				secondary += "S"
				current += 2
			} else if current+1 < length && word[current+1] == 'H' {
				// 'CH' -> 'X' (как в 'chocolate')
				primary += "X"
				secondary += "X"
				current += 2
			} else {
				// Другие случаи 'C' -> 'K'
				primary += "K"
				secondary += "K"
				current++
			}

		case 'D':
			if current+2 < length && word[current+1] == 'G' &&
				(word[current+2] == 'E' || word[current+2] == 'I' || word[current+2] == 'Y') {
				// 'DGE', 'DGI', 'DGY' -> 'J'
				primary += "J"
				secondary += "J"
				current += 3
			} else {
				// Другие случаи 'D' -> 'T'
				primary += "T"
				secondary += "T"
				current++
			}

		case 'F':
			primary += "F"
			secondary += "F"
			current++

		case 'G':
			if current+1 < length && word[current+1] == 'H' {
				if current > 0 && !isVowel(word[current-1]) {
					// 'GH' не в начале слова и после согласной -> 'K'
					primary += "K"
					secondary += "K"
				} else if current == 0 {
					// 'GH' в начале слова -> 'K'
					if current+2 >= length || word[current+2] != 'I' {
						primary += "K"
						secondary += "K"
					}
				}
				current += 2
			} else if current+1 < length && word[current+1] == 'N' {
				if current == 0 && isVowel(word[current+2]) {
					// 'GN' в начале слова перед гласной -> 'N'
					primary += "KN"
					secondary += "N"
				} else {
					// Другие случаи 'GN' -> 'N'
					primary += "N"
					secondary += "N"
				}
				current += 2
			} else if current+1 < length &&
				(word[current+1] == 'E' || word[current+1] == 'I' || word[current+1] == 'Y') {
				// 'GE', 'GI', 'GY' -> 'J'
				primary += "J"
				secondary += "J"
				current += 2
			} else {
				// Другие случаи 'G' -> 'K'
				primary += "K"
				secondary += "K"
				current++
			}

		case 'H':
			// 'H' сохраняется только если после гласной и не перед гласной
			if current > 0 && isVowel(word[current-1]) &&
				(current+1 >= length || !isVowel(word[current+1])) {
				primary += "H"
				secondary += "H"
			}
			current++

		case 'J':
			// 'J' -> 'J' (звучит как 'Y' в некоторых языках)
			primary += "J"
			secondary += "J"
			current++

		case 'K':
			// 'K' -> 'K'
			if current > 0 && word[current-1] == 'C' {
				// Пропускаем 'K' после 'C'
				current++
			} else {
				primary += "K"
				secondary += "K"
				current++
			}

		case 'L':
			primary += "L"
			secondary += "L"
			current++

		case 'M':
			primary += "M"
			secondary += "M"
			current++

		case 'N':
			primary += "N"
			secondary += "N"
			current++

		case 'P':
			if current+1 < length && word[current+1] == 'H' {
				// 'PH' -> 'F'
				primary += "F"
				secondary += "F"
				current += 2
			} else {
				primary += "P"
				secondary += "P"
				current++
			}

		case 'Q':
			primary += "K"
			secondary += "K"
			current++

		case 'R':
			primary += "R"
			secondary += "R"
			current++

		case 'S':
			if current+2 < length && word[current+1] == 'I' &&
				(word[current+2] == 'O' || word[current+2] == 'A') {
				// 'SIO', 'SIA' -> 'X'
				primary += "X"
				secondary += "S"
				current += 3
			} else if current+1 < length && word[current+1] == 'H' {
				// 'SH' -> 'X'
				primary += "X"
				secondary += "X"
				current += 2
			} else {
				primary += "S"
				secondary += "S"
				current++
			}

		case 'T':
			if current+2 < length && word[current+1] == 'I' &&
				(word[current+2] == 'O' || word[current+2] == 'A') {
				// 'TIO', 'TIA' -> 'X'
				primary += "X"
				secondary += "X"
				current += 3
			} else if current+1 < length && word[current+1] == 'H' {
				// 'TH' -> '0' (специальный код для 'th')
				primary += "0"
				secondary += "0"
				current += 2
			} else {
				primary += "T"
				secondary += "T"
				current++
			}

		case 'V':
			primary += "F"
			secondary += "F"
			current++

		case 'W':
			// 'W' сохраняется только если перед гласной
			if current+1 < length && isVowel(word[current+1]) {
				primary += "W"
				secondary += "W"
			}
			current++

		case 'X':
			// 'X' -> 'KS'
			primary += "KS"
			secondary += "KS"
			current++

		case 'Z':
			primary += "S"
			secondary += "S"
			current++

		default:
			current++
		}
	}

	// Ограничиваем длину кодов
	if len(primary) > 4 {
		primary = primary[:4]
	}
	if len(secondary) > 4 {
		secondary = secondary[:4]
	}

	// Если вторичный код идентичен первичному, возвращаем пустую строку
	if primary == secondary {
		secondary = ""
	}

	return primary, secondary
}

// cleanString удаляет не-буквенные символы и приводит к верхнему регистру
func cleanString(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			result.WriteRune(unicode.ToUpper(r))
		}
	}
	return result.String()
}

// isVowel проверяет, является ли символ гласной
func isVowel(c byte) bool {
	return c == 'A' || c == 'E' || c == 'I' || c == 'O' || c == 'U' || c == 'Y'
}
