package translit

import (
	"strings"
	"unicode"
)

// Проверяет, содержит ли строка кириллические символы (русские или украинские)
func IsCyrillic(text string) bool {
	// Карта украинских символов, которые могут не входить в стандартный диапазон unicode.Cyrillic
	ukrainianChars := map[rune]bool{
		'і': true, 'ї': true, 'є': true, 'ґ': true,
		'І': true, 'Ї': true, 'Є': true, 'Ґ': true,
	}

	for _, r := range text {
		if unicode.Is(unicode.Cyrillic, r) {
			return true
		}
		// Проверка украинских символов
		if ukrainianChars[r] {
			return true
		}
	}
	return false
}

// GetAllTransliterations возвращает все варианты транслитерации для имени
func GetAllTransliterations(name string) []string {
	// Если имя не на кириллице, возвращаем его как есть
	if !IsCyrillic(name) {
		return []string{name}
	}

	// Получаем все стандартные транслитерации
	standardTransliterations := []string{
		TranslitGOST(name),
		TranslitISO9(name),
		TranslitBGNPCGN(name),
		TranslitUNGEGN(name),
	}

	// Добавляем исходное имя в список
	allTransliterations := append([]string{name}, standardTransliterations...)

	// Убираем дубликаты
	uniqueTransliterations := removeDuplicates(allTransliterations)

	return uniqueTransliterations
}

// removeDuplicates удаляет дубликаты из слайса строк, сохраняя порядок
func removeDuplicates(elements []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, element := range elements {
		if !seen[element] {
			seen[element] = true
			result = append(result, element)
		}
	}
	return result
}

// GetTranslitFunction возвращает функцию транслитерации по имени стандарта
func GetTranslitFunction(standard string) func(string) string {
	switch standard {
	case "iso9":
		return TranslitISO9
	case "gost":
		return TranslitGOST
	case "bgnpcgn":
		return TranslitBGNPCGN
	case "ungegn":
		return TranslitUNGEGN
	case "ukrainian":
		return TranslitUkrainian
	default:
		return TranslitISO9 // По умолчанию ISO 9
	}
}

// TranslitUkrainian транслитерация по украинскому стандарту
func TranslitUkrainian(text string) string {
	// Приводим к нижнему регистру
	text = strings.ToLower(text)

	mapping := map[rune]string{
		// Общие с русским буквы
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ж': "zh", 'з': "z", 'и': "y", 'к': "k", 'л': "l", 'м': "m",
		'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
		'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh",
		'щ': "shch", 'ь': "", 'ю': "yu", 'я': "ya",

		// Специфические украинские буквы
		'і': "i", 'ї': "yi", 'є': "ye", 'ґ': "g",
	}

	mapping['й'] = "y" // Изменено на "y" для соответствия украинской транслитерации

	// Специальные правила для украинской транслитерации
	// Обработка 'зг' как 'zgh'
	text = strings.ReplaceAll(text, "зг", "zgh")

	return transliterate(text, mapping)
}

// transliterate вспомогательная функция для транслитерации
func transliterate(text string, mapping map[rune]string) string {
	var sb strings.Builder
	for _, r := range text {
		if translit, ok := mapping[r]; ok {
			sb.WriteString(translit)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// GetAllTransliterationsWithTypos возвращает все варианты транслитерации для имени,
// включая возможные опечатки
func GetAllTransliterationsWithTypos(name string) []string {
	// Получаем все стандартные транслитерации
	variations := GetAllTransliterations(name)

	// Добавляем вариации с типичными опечатками
	result := make([]string, 0, len(variations)*3)
	result = append(result, variations...)

	// Типичные опечатки в транслитерации
	for _, v := range variations {
		// Замена 'y' на 'i' и наоборот
		result = append(result, strings.ReplaceAll(v, "y", "i"))
		result = append(result, strings.ReplaceAll(v, "i", "y"))

		// Замена 'e' на 'i' и наоборот (для имен типа Sergey/Sergei)
		if strings.HasSuffix(v, "ey") {
			result = append(result, strings.TrimSuffix(v, "ey")+"ei")
		}
		if strings.HasSuffix(v, "ei") {
			result = append(result, strings.TrimSuffix(v, "ei")+"ey")
		}

		// Замена 'sh' на 'sch' и наоборот
		result = append(result, strings.ReplaceAll(v, "sh", "sch"))
		result = append(result, strings.ReplaceAll(v, "sch", "sh"))

		// Замена 'zh' на 'j' и наоборот
		result = append(result, strings.ReplaceAll(v, "zh", "j"))
		result = append(result, strings.ReplaceAll(v, "j", "zh"))

		// Замена 'ts' на 'c' и наоборот
		result = append(result, strings.ReplaceAll(v, "ts", "c"))
		result = append(result, strings.ReplaceAll(v, "c", "ts"))

		// Замена 'kh' на 'h' и наоборот
		result = append(result, strings.ReplaceAll(v, "kh", "h"))
		result = append(result, strings.ReplaceAll(v, "h", "kh"))

		// Замена 'yu' на 'iu' и наоборот
		result = append(result, strings.ReplaceAll(v, "yu", "iu"))
		result = append(result, strings.ReplaceAll(v, "iu", "yu"))

		// Замена 'ya' на 'ia' и наоборот
		result = append(result, strings.ReplaceAll(v, "ya", "ia"))
		result = append(result, strings.ReplaceAll(v, "ia", "ya"))

		// Удвоение согласных (типичная ошибка)
		for _, c := range []string{"n", "l", "t", "s", "r", "p"} {
			result = append(result, strings.ReplaceAll(v, c, c+c))
		}

		// Пропуск гласных (типичная ошибка)
		for _, c := range []string{"a", "e", "i", "o", "u", "y"} {
			result = append(result, strings.ReplaceAll(v, c, ""))
		}
	}

	// Удаляем дубликаты
	uniqueVariations := make(map[string]bool)
	for _, v := range result {
		if v != "" {
			uniqueVariations[v] = true
		}
	}

	// Преобразуем обратно в слайс
	finalResult := make([]string, 0, len(uniqueVariations))
	for v := range uniqueVariations {
		finalResult = append(finalResult, v)
	}

	return finalResult
}
