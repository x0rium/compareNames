package translit

import (
	"strings"
)

// TranslitUNGEGN транслитерация по системе UNGEGN
func TranslitUNGEGN(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	mapping := map[rune]string{
		// Русские буквы
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ё': "ë", 'ж': "ž", 'з': "z", 'и': "i", 'й': "j", 'к': "k",
		'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "h", 'ц': "c",
		'ч': "č", 'ш': "š", 'щ': "šč", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "è", 'ю': "ju", 'я': "ja",

		// Украинские буквы
		'і': "i", 'ї': "ji", 'є': "je", 'ґ': "g",
	}

	return transliterate(text, mapping)
}

// TranslitUNGEGNSimplified упрощенная транслитерация по UNGEGN без диакритических знаков
func TranslitUNGEGNSimplified(text string) string {
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

// TranslitUNGEGNReverse обратная транслитерация с латиницы на кириллицу по UNGEGN
func TranslitUNGEGNReverse(text string) string {
	// Сначала заменяем многосимвольные комбинации (для упрощенной версии)
	replacements := map[string]string{
		"shch": "щ", "shh": "щ", "zh": "ж", "ch": "ч", "sh": "ш",
		"yu": "ю", "ya": "я", "ye": "є", "yi": "ї",
		"ju": "ю", "ja": "я", "je": "є", "ji": "ї",
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
		'č': "ч", 'š': "ш", 'è': "э", 'y': "ы",
		'ž': "ж",
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

// GetUNGEGNVariations возвращает возможные вариации транслитерации по UNGEGN
func GetUNGEGNVariations(text string) []string {
	if !IsCyrillic(text) {
		// Если текст не на кириллице, пытаемся выполнить обратную транслитерацию
		return []string{TranslitUNGEGNReverse(text)}
	}

	// Основная транслитерация
	base := TranslitUNGEGN(text)
	simplified := TranslitUNGEGNSimplified(text)

	variations := []string{base, simplified}

	// Вариации для часто встречающихся случаев
	variations = append(variations, strings.ReplaceAll(simplified, "zh", "j"))
	variations = append(variations, strings.ReplaceAll(simplified, "ch", "tch"))
	variations = append(variations, strings.ReplaceAll(simplified, "sh", "sch"))
	variations = append(variations, strings.ReplaceAll(simplified, "yu", "iu"))
	variations = append(variations, strings.ReplaceAll(simplified, "ya", "ia"))

	return variations
}
