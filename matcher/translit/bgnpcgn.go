package translit

import (
	"strings"
)

// TranslitBGNPCGN транслитерация по системе BGN/PCGN
func TranslitBGNPCGN(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	mapping := map[rune]string{
		// Русские буквы
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ё': "ë", 'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k",
		'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
		'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "'",
		'э': "e", 'ю': "yu", 'я': "ya",

		// Украинские буквы
		'і': "i", 'ї': "yi", 'є': "ye", 'ґ': "g",
	}

	// Специальная обработка для 'е' в начале слова, после гласных и после ь, ъ
	words := strings.Fields(text)
	var result []string

	for _, word := range words {
		var sb strings.Builder
		for i, r := range word {
			if r == 'е' {
				// 'е' в начале слова
				if i == 0 {
					sb.WriteString("ye")
					continue
				}

				// 'е' после гласных или ь, ъ
				prevChar := rune(word[i-1])
				if isVowel(prevChar) || prevChar == 'ь' || prevChar == 'ъ' {
					sb.WriteString("ye")
					continue
				}
			}

			if translit, ok := mapping[r]; ok {
				sb.WriteString(translit)
			} else {
				sb.WriteRune(r)
			}
		}
		result = append(result, sb.String())
	}

	return strings.Join(result, " ")
}

// TranslitBGNPCGNReverse обратная транслитерация с латиницы на кириллицу по BGN/PCGN
func TranslitBGNPCGNReverse(text string) string {
	// Сначала заменяем многосимвольные комбинации
	replacements := map[string]string{
		"shch": "щ", "zh": "ж", "kh": "х", "ts": "ц",
		"ch": "ч", "sh": "ш", "ye": "е", "yu": "ю", "ya": "я",
		"yi": "ї",
	}

	for latin, cyrillic := range replacements {
		text = strings.ReplaceAll(text, latin, cyrillic)
	}

	// Затем заменяем одиночные символы
	mapping := map[rune]string{
		'a': "а", 'b': "б", 'v': "в", 'g': "г", 'd': "д", 'e': "е",
		'ë': "ё", 'z': "з", 'i': "и", 'y': "й", 'k': "к", 'l': "л",
		'm': "м", 'n': "н", 'o': "о", 'p': "п", 'r': "р", 's': "с",
		't': "т", 'u': "у", 'f': "ф", '\'': "ь",
	}

	var sb strings.Builder
	for _, r := range text {
		if cyrillic, ok := mapping[r]; ok {
			sb.WriteString(cyrillic)
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// GetBGNPCGNVariations возвращает возможные вариации транслитерации по BGN/PCGN
func GetBGNPCGNVariations(text string) []string {
	if !IsCyrillic(text) {
		// Если текст не на кириллице, пытаемся выполнить обратную транслитерацию
		return []string{TranslitBGNPCGNReverse(text)}
	}

	// Основная транслитерация
	base := TranslitBGNPCGN(text)
	variations := []string{base}

	// Вариации для часто встречающихся случаев
	variations = append(variations, strings.ReplaceAll(base, "y", "j"))
	variations = append(variations, strings.ReplaceAll(base, "kh", "h"))
	variations = append(variations, strings.ReplaceAll(base, "'", ""))
	variations = append(variations, strings.ReplaceAll(base, "yu", "iu"))
	variations = append(variations, strings.ReplaceAll(base, "ya", "ia"))
	variations = append(variations, strings.ReplaceAll(base, "ye", "ie"))

	// Вариации для окончаний
	if strings.HasSuffix(base, "y") {
		variations = append(variations, strings.TrimSuffix(base, "y")+"i")
		variations = append(variations, strings.TrimSuffix(base, "y")+"iy")
	}

	return variations
}

// isVowel проверяет, является ли символ гласной
func isVowel(r rune) bool {
	vowels := []rune{'а', 'е', 'ё', 'и', 'о', 'у', 'ы', 'э', 'ю', 'я', 'і', 'ї', 'є'}
	for _, v := range vowels {
		if r == v {
			return true
		}
	}
	return false
}
