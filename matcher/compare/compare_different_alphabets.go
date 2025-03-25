package compare

import (
	"strings"
	"time"

	"github.com/x0rium/compareNames/matcher/similarity"
	"github.com/x0rium/compareNames/matcher/translit"
)

// DifferentAlphabetsResult содержит результат сравнения имен на разных алфавитах
type DifferentAlphabetsResult struct {
	ExactMatch           bool    // Точное совпадение
	Score                int     // Оценка совпадения (0-100)
	MatchType            string  // Тип совпадения ("match", "possible_match", "no_match")
	BestMatch1           string  // Лучшее совпадение для первого имени
	BestMatch2           string  // Лучшее совпадение для второго имени
	LevenshteinScore     float64 // Оценка по Левенштейну
	JaroWinklerScore     float64 // Оценка по Джаро-Винклеру
	PhoneticScore        float64 // Фонетическая оценка
	DoubleMetaphoneScore float64 // Оценка по Double Metaphone
	ProcessingTimeMS     int64   // Время обработки в миллисекундах
}

// CompareDifferentAlphabets сравнивает имена на разных алфавитах (кириллица и латиница)
func CompareDifferentAlphabets(name1, name2 string, isName1Cyrillic bool, startTime time.Time,
	getNameVariations func(string) []string, normalizeNameParts func(string) []string,
	compareNameParts func([]string, []string) float64, minExactMatchScore int) DifferentAlphabetsResult {

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

	// Для каждой вариации кириллического имени проверяем все транслитерации
	for _, cyrillicVar := range cyrillicVariations {
		// Получаем все возможные транслитерации
		transliterations := translit.GetAllTransliterationsWithTypos(cyrillicVar)

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

					return DifferentAlphabetsResult{
						ExactMatch:           false,
						Score:                95, // Высокий балл для точных транслитераций
						MatchType:            "match",
						BestMatch1:           originalCyrillic,
						BestMatch2:           originalLatin,
						LevenshteinScore:     0.95,
						JaroWinklerScore:     0.95,
						PhoneticScore:        1.0,
						DoubleMetaphoneScore: 1.0,
						ProcessingTimeMS:     time.Since(startTime).Milliseconds(),
					}
				}

				// Проверяем на частичное совпадение (нечеткое)
				// Разбиваем имена на части и сравниваем каждую часть
				transParts := normalizeNameParts(trans)
				latinParts := normalizeNameParts(latinVar)

				// Если структура имен слишком разная, продолжаем
				if abs(len(transParts)-len(latinParts)) > 1 {
					continue
				}

				// Вычисляем схожесть между частями имен
				partSimilarity := compareNameParts(transParts, latinParts)

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
				}
			}
		}
	}

	// Если нашли достаточно высокое совпадение по транслитерации
	if maxSimilarity > 0.7 {
		score := int(maxSimilarity * 100)
		matchType := "match"
		if score < minExactMatchScore {
			matchType = "possible_match"
		}

		// Вычисляем отдельные оценки для лучшего совпадения
		levScore := similarity.LevenshteinSimilarity(bestTransliteration, bestLatin)
		jaroScore := similarity.JaroWinklerSimilarity(bestTransliteration, bestLatin)
		phoneticScore := similarity.PhoneticSimilarity(bestTransliteration, bestLatin)
		doubleMetaphoneScore := similarity.DoubleMetaphoneSimilarity(bestTransliteration, bestLatin)

		return DifferentAlphabetsResult{
			ExactMatch:           false,
			Score:                score,
			MatchType:            matchType,
			BestMatch1:           bestCyrillic,
			BestMatch2:           bestLatin,
			LevenshteinScore:     levScore,
			JaroWinklerScore:     jaroScore,
			PhoneticScore:        phoneticScore,
			DoubleMetaphoneScore: doubleMetaphoneScore,
			ProcessingTimeMS:     time.Since(startTime).Milliseconds(),
		}
	}

	// Если не нашли хорошего совпадения, возвращаем пустой результат
	return DifferentAlphabetsResult{
		ExactMatch:       false,
		Score:            0,
		MatchType:        "no_match",
		ProcessingTimeMS: time.Since(startTime).Milliseconds(),
	}
}

// abs вычисляет абсолютное значение разницы
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
