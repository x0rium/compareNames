package compare

import (
	"strings"
	"time"

	"github.com/x0rium/compareNames/matcher/similarity"
	"github.com/x0rium/compareNames/matcher/translit"
)

// TransliterationResult содержит результат сравнения имен на разных алфавитах
type TransliterationResult struct {
	ExactMatch           bool    // Точное совпадение
	Score                int     // Оценка совпадения (0-100)
	MatchType            string  // Тип совпадения ("match", "possible_match", "no_match")
	BestMatch1           string  // Лучшее совпадение для первого имени
	BestMatch2           string  // Лучшее совпадение для второго имени
	BestTransliteration  string  // Лучшая транслитерация
	LevenshteinScore     float64 // Оценка по алгоритму Левенштейна
	JaroWinklerScore     float64 // Оценка по алгоритму Джаро-Винклера
	PhoneticScore        float64 // Фонетическая оценка
	DoubleMetaphoneScore float64 // Оценка по алгоритму Double Metaphone
	ProcessingTime       int64   // Время обработки в миллисекундах
}

// CompareTransliteration сравнивает имена на разных алфавитах (кириллица и латиница)
func CompareTransliteration(name1, name2 string, isName1Cyrillic bool, startTime time.Time) TransliterationResult {
	// Определяем кириллическое и латинское имя
	cyrillicName := name1
	latinName := name2
	if !isName1Cyrillic {
		cyrillicName = name2
		latinName = name1
	}

	// Получаем вариации имен (с перестановками)
	cyrillicVariations := getNameVariations(cyrillicName)
	latinVariations := getNameVariations(latinName)

	// Проверяем все варианты транслитераций и их совпадения
	maxSimilarity := 0.0
	bestCyrillic := ""
	bestLatin := ""
	bestTransliteration := ""
	var bestLevScore, bestJaroScore, bestPhoneticScore, bestDMScore float64

	// Приоритет для полных имен (с тремя частями)
	cyrillicFullNames := make([]string, 0)
	latinFullNames := make([]string, 0)

	// Отбираем полные имена (с тремя частями)
	for _, cyrillicVar := range cyrillicVariations {
		parts := strings.Fields(cyrillicVar)
		if len(parts) == 3 {
			cyrillicFullNames = append(cyrillicFullNames, cyrillicVar)
		}
	}

	for _, latinVar := range latinVariations {
		parts := strings.Fields(latinVar)
		if len(parts) == 3 {
			latinFullNames = append(latinFullNames, latinVar)
		}
	}

	// Если есть полные имена, используем только их
	if len(cyrillicFullNames) > 0 && len(latinFullNames) > 0 {
		cyrillicVariations = cyrillicFullNames
		latinVariations = latinFullNames
	}

	// Ограничиваем количество вариаций для улучшения производительности
	maxVariations := 5
	if len(cyrillicVariations) > maxVariations {
		cyrillicVariations = cyrillicVariations[:maxVariations]
	}
	if len(latinVariations) > maxVariations {
		latinVariations = latinVariations[:maxVariations]
	}

	// Для каждой вариации кириллического имени проверяем все транслитерации
	for _, cyrillicVar := range cyrillicVariations {
		// Получаем все возможные транслитерации, включая опечатки
		transliterations := translit.GetAllTransliterationsWithTypos(cyrillicVar)

		// Ограничиваем количество транслитераций для улучшения производительности
		maxTransliterations := 5
		if len(transliterations) > maxTransliterations {
			transliterations = transliterations[:maxTransliterations]
		}

		// Проверяем все транслитерации против всех вариаций латинского имени
		for _, trans := range transliterations {
			for _, latinVar := range latinVariations {
				// Проверяем на точное совпадение (с учетом регистра)
				if strings.ToLower(trans) == strings.ToLower(latinVar) {
					// Используем исходные имена для результата, если это полные имена
					originalCyrillic := cyrillicVar
					originalLatin := latinVar

					// Если это не полные имена, но есть полные имена в вариациях,
					// используем первое полное имя из вариаций
					if len(strings.Fields(cyrillicVar)) < 3 && len(cyrillicFullNames) > 0 {
						originalCyrillic = cyrillicFullNames[0]
					}
					if len(strings.Fields(latinVar)) < 3 && len(latinFullNames) > 0 {
						originalLatin = latinFullNames[0]
					}

					return TransliterationResult{
						ExactMatch:           false,
						Score:                95, // Высокий балл для точных транслитераций
						MatchType:            "match",
						BestMatch1:           originalCyrillic,
						BestMatch2:           originalLatin,
						BestTransliteration:  trans,
						LevenshteinScore:     0.95,
						JaroWinklerScore:     0.95,
						PhoneticScore:        1.0,
						DoubleMetaphoneScore: 1.0,
						ProcessingTime:       time.Since(startTime).Milliseconds(),
					}
				}

				// Проверяем на частичное совпадение (нечеткое)
				// Разбиваем имена на части и сравниваем каждую часть
				transParts := normalizeNameParts(trans)
				latinParts := normalizeNameParts(latinVar)

				// Если структура имен слишком разная, продолжаем
				if absInt(len(transParts)-len(latinParts)) > 1 {
					continue
				}

				// Вычисляем схожесть между частями имен
				partSimilarity := CompareNameParts(transParts, latinParts, nil)

				// Рассчитываем схожесть между строками целиком для дополнительной проверки
				levScore := similarity.LevenshteinSimilarity(trans, latinVar)
				jaroScore := similarity.JaroWinklerSimilarity(trans, latinVar)
				phoneticScore := similarity.PhoneticSimilarity(trans, latinVar)
				doubleMetaphoneScore := similarity.DoubleMetaphoneSimilarity(trans, latinVar)

				// Вычисляем схожесть между строками
				stringSimilarity := (levScore*0.3 + jaroScore*0.3 +
					phoneticScore*0.2 + doubleMetaphoneScore*0.2)

				// Вычисляем итоговую схожесть, отдавая приоритет сравнению по частям
				totalSimilarity := partSimilarity*0.6 + stringSimilarity*0.4

				// Для скоринговой системы ФЗ 115 снижаем оценку для неточных совпадений
				if totalSimilarity < 0.9 {
					// Применяем штраф в зависимости от уровня схожести
					if totalSimilarity > 0.8 {
						totalSimilarity *= 0.95 // Небольшой штраф для высокой схожести
					} else if totalSimilarity > 0.7 {
						totalSimilarity *= 0.9 // Средний штраф
					} else {
						totalSimilarity *= 0.8 // Значительный штраф для низкой схожести
					}
				}

				if totalSimilarity > maxSimilarity {
					maxSimilarity = totalSimilarity
					bestCyrillic = cyrillicVar
					bestLatin = latinVar
					bestTransliteration = trans
					bestLevScore = levScore
					bestJaroScore = jaroScore
					bestPhoneticScore = phoneticScore
					bestDMScore = doubleMetaphoneScore
				}
			}
		}
	}

	// Если нашли достаточно высокое совпадение по транслитерации
	if maxSimilarity > 0.7 {
		score := int(maxSimilarity * 100)
		matchType := "match"
		if score < 90 { // Используем константу MinExactMatchScore
			matchType = "possible_match"
		}

		return TransliterationResult{
			ExactMatch:           false,
			Score:                score,
			MatchType:            matchType,
			BestMatch1:           bestCyrillic,
			BestMatch2:           bestLatin,
			BestTransliteration:  bestTransliteration,
			LevenshteinScore:     bestLevScore,
			JaroWinklerScore:     bestJaroScore,
			PhoneticScore:        bestPhoneticScore,
			DoubleMetaphoneScore: bestDMScore,
			ProcessingTime:       time.Since(startTime).Milliseconds(),
		}
	}

	// Если не нашли хорошего совпадения
	return TransliterationResult{
		ExactMatch:     false,
		Score:          int(maxSimilarity * 100),
		MatchType:      "no_match",
		ProcessingTime: time.Since(startTime).Milliseconds(),
	}
}

// getNameVariations генерирует различные вариации имени, включая перестановки
func getNameVariations(name string) []string {
	parts := normalizeNameParts(name)
	if len(parts) == 0 {
		return []string{}
	}

	variationsMap := make(map[string]bool) // Используем map для удаления дубликатов

	// Добавляем исходное имя
	original := strings.Join(parts, " ")
	variationsMap[original] = true

	// Проверяем на наличие инициалов
	hasInitials := HasInitials(name)

	// Если есть хотя бы 2 части, добавляем перестановки
	if len(parts) >= 2 {
		// Добавляем все возможные перестановки для 2 и 3 частей
		if len(parts) == 2 {
			// ИФ или ФИ
			variationsMap[strings.Join([]string{parts[1], parts[0]}, " ")] = true
		} else if len(parts) == 3 {
			// ФИО -> ИФО
			variationsMap[strings.Join([]string{parts[1], parts[0], parts[2]}, " ")] = true
			// ФИО -> ФОИ
			variationsMap[strings.Join([]string{parts[0], parts[2], parts[1]}, " ")] = true
			// ФИО -> ОФИ
			variationsMap[strings.Join([]string{parts[2], parts[0], parts[1]}, " ")] = true
			// ФИО -> ОИФ
			variationsMap[strings.Join([]string{parts[2], parts[1], parts[0]}, " ")] = true
			// ФИО -> ИОФ
			variationsMap[strings.Join([]string{parts[1], parts[2], parts[0]}, " ")] = true

			// Также добавляем вариации с пропущенными частями
			// Без отчества
			variationsMap[strings.Join([]string{parts[0], parts[1]}, " ")] = true
			variationsMap[strings.Join([]string{parts[1], parts[0]}, " ")] = true
			// Без имени
			variationsMap[strings.Join([]string{parts[0], parts[2]}, " ")] = true
			variationsMap[strings.Join([]string{parts[2], parts[0]}, " ")] = true
			// Без фамилии
			variationsMap[strings.Join([]string{parts[1], parts[2]}, " ")] = true
			variationsMap[strings.Join([]string{parts[2], parts[1]}, " ")] = true
		}
	}

	// Если есть инициалы, добавляем специальную обработку
	if hasInitials {
		// Определяем, какие части являются инициалами
		var fullParts, initialParts []string
		for _, p := range parts {
			if strings.HasSuffix(p, ".") || len(p) == 1 {
				initialParts = append(initialParts, p)
			} else {
				fullParts = append(fullParts, p)
			}
		}

		// Комбинируем их в разных порядках
		if len(fullParts) > 0 && len(initialParts) > 0 {
			// Сначала полные части, потом инициалы
			variationsMap[strings.Join(append(fullParts, initialParts...), " ")] = true
			// Сначала инициалы, потом полные части
			variationsMap[strings.Join(append(initialParts, fullParts...), " ")] = true

			// Для каждой полной части добавляем вариацию с инициалом
			for _, fullPart := range fullParts {
				for _, initial := range initialParts {
					variationsMap[fullPart+" "+initial] = true
					variationsMap[initial+" "+fullPart] = true
				}
			}
		}
	}

	// Преобразуем map в слайс
	variations := make([]string, 0, len(variationsMap))
	for v := range variationsMap {
		variations = append(variations, v)
	}

	return variations
}

// normalizeNameParts разбивает имя на части и возвращает их в нормализованном виде
func normalizeNameParts(name string) []string {
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

// absInt возвращает абсолютное значение числа
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
