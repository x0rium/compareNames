package translit

import (
	"strings"
)

// TranslitISO9 транслитерация по стандарту ISO 9
func TranslitISO9(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	mapping := map[rune]string{
		// Русские буквы
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ё': "ë", 'ж': "ž", 'з': "z", 'и': "i", 'й': "j", 'к': "k",
		'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "h", 'ц': "c",
		'ч': "č", 'ш': "š", 'щ': "ŝ", 'ъ': "ʺ", 'ы': "y", 'ь': "ʹ",
		'э': "è", 'ю': "û", 'я': "â",

		// Украинские буквы
		'і': "i", 'ї': "ï", 'є': "ê", 'ґ': "g",
	}

	return transliterate(text, mapping)
}

// TranslitISO9Simplified упрощенная транслитерация по ISO 9 без диакритических знаков
func TranslitISO9Simplified(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	mapping := map[rune]string{
		// Русские буквы
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ё': "e", 'ж': "zh", 'з': "z", 'и': "i", 'й': "j", 'к': "k",
		'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "h", 'ц': "c",
		'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",

		// Украинские буквы
		'і': "i", 'ї': "yi", 'є': "ye", 'ґ': "g",
	}

	return transliterate(text, mapping)
}

// TranslitISO9Reverse обратная транслитерация с латиницы на кириллицу по ISO 9
func TranslitISO9Reverse(text string) string {
	// Сначала заменяем многосимвольные комбинации (для упрощенной версии)
	replacements := map[string]string{
		"shch": "щ", "zh": "ж", "ch": "ч", "sh": "ш",
		"yu": "ю", "ya": "я", "ye": "є", "yi": "ї",
	}

	for latin, cyrillic := range replacements {
		text = strings.ReplaceAll(text, latin, cyrillic)
	}

	// Затем заменяем одиночные символы и символы с диакритическими знаками
	mapping := map[rune]string{
		'a': "а", 'b': "б", 'v': "в", 'g': "г", 'd': "д", 'e': "е",
		'ë': "ё", 'z': "з", 'i': "и", 'j': "й", 'k': "к", 'l': "л",
		'm': "м", 'n': "н", 'o': "о", 'p': "п", 'r': "р", 's': "с",
		't': "т", 'u': "у", 'f': "ф", 'h': "х", 'c': "ц",
		'č': "ч", 'š': "ш", 'ŝ': "щ", 'ʺ': "ъ", 'y': "ы", 'ʹ': "ь",
		'è': "э", 'û': "ю", 'â': "я",
		'ê': "є", 'ï': "ї",
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

// GetISO9Variations возвращает возможные вариации транслитерации по ISO 9
func GetISO9Variations(text string) []string {
	if !IsCyrillic(text) {
		// Если текст не на кириллице, пытаемся выполнить обратную транслитерацию
		return []string{TranslitISO9Reverse(text)}
	}

	// Основная транслитерация
	base := TranslitISO9(text)
	simplified := TranslitISO9Simplified(text)

	variations := []string{base, simplified}

	// Вариации для часто встречающихся случаев
	variations = append(variations, strings.ReplaceAll(simplified, "zh", "j"))
	variations = append(variations, strings.ReplaceAll(simplified, "ch", "tch"))
	variations = append(variations, strings.ReplaceAll(simplified, "sh", "sch"))
	variations = append(variations, strings.ReplaceAll(simplified, "yu", "iu"))
	variations = append(variations, strings.ReplaceAll(simplified, "ya", "ia"))

	return variations
}
