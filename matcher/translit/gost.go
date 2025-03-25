package translit

import (
	"strings"
)

// TranslitGOST транслитерация по ГОСТ 7.79-2000 (система Б)
func TranslitGOST(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	// Карта соответствия символов
	mapping := map[rune]string{
		// Русские буквы
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i", 'й': "j", 'к': "k",
		'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "cz",
		'ч': "ch", 'ш': "sh", 'щ': "shh", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",

		// Украинские буквы
		'є': "ye", 'і': "i", 'ї': "yi", 'ґ': "g",
	}

	// Специальные случаи для окончаний
	if strings.HasSuffix(text, "ий") {
		text = strings.TrimSuffix(text, "ий") + "y"
	}
	if strings.HasSuffix(text, "ый") {
		text = strings.TrimSuffix(text, "ый") + "y"
	}

	// Применяем транслитерацию
	var sb strings.Builder
	for _, r := range text {
		if translit, ok := mapping[r]; ok {
			sb.WriteString(translit)
		} else {
			sb.WriteRune(r)
		}
	}

	// Постобработка для специфических случаев
	result := sb.String()

	// Обработка сочетаний букв
	result = strings.ReplaceAll(result, "cz", "ts")    // Заменяем "cz" на "ts" для лучшей читаемости
	result = strings.ReplaceAll(result, "shh", "shch") // Заменяем "shh" на "shch" для лучшей читаемости

	// Обработка сочетаний с 'е' после согласных
	consonants := []string{"b", "v", "g", "d", "z", "k", "l", "m", "n", "p", "r", "s", "t", "f", "kh"}
	for _, consonant := range consonants {
		result = strings.ReplaceAll(result, consonant+"e", consonant+"e")
	}

	return result
}

// TranslitGOSTReverse обратная транслитерация с латиницы на кириллицу по ГОСТ 7.79-2000
func TranslitGOSTReverse(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	// Сначала заменяем многосимвольные комбинации
	replacements := map[string]string{
		"shch": "щ", "shh": "щ",
		"zh": "ж", "kh": "х", "ts": "ц", "cz": "ц",
		"ch": "ч", "sh": "ш",
		"yu": "ю", "ya": "я", "yo": "ё",
		"ye": "е", "yi": "ї",
	}

	for latin, cyrillic := range replacements {
		text = strings.ReplaceAll(text, latin, cyrillic)
	}

	// Затем заменяем одиночные символы
	mapping := map[rune]string{
		'a': "а", 'b': "б", 'v': "в", 'g': "г", 'd': "д", 'e': "е",
		'z': "з", 'i': "и", 'j': "й", 'k': "к", 'l': "л", 'm': "м",
		'n': "н", 'o': "о", 'p': "п", 'r': "р", 's': "с", 't': "т",
		'u': "у", 'f': "ф", 'y': "ы",
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

// GetGOSTVariations возвращает возможные вариации транслитерации по ГОСТ
func GetGOSTVariations(text string) []string {
	if !IsCyrillic(text) {
		// Если текст не на кириллице, пытаемся выполнить обратную транслитерацию
		return []string{TranslitGOSTReverse(text)}
	}

	// Основная транслитерация
	base := TranslitGOST(text)
	variations := []string{base}

	// Вариации для часто встречающихся случаев
	variations = append(variations, strings.ReplaceAll(base, "ts", "c"))
	variations = append(variations, strings.ReplaceAll(base, "kh", "h"))
	variations = append(variations, strings.ReplaceAll(base, "yo", "e"))
	variations = append(variations, strings.ReplaceAll(base, "yu", "iu"))
	variations = append(variations, strings.ReplaceAll(base, "ya", "ia"))

	// Вариации для окончаний
	if strings.HasSuffix(base, "y") {
		variations = append(variations, strings.TrimSuffix(base, "y")+"iy")
		variations = append(variations, strings.TrimSuffix(base, "y")+"yy")
	}

	return variations
}
